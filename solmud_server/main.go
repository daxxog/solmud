package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// ---------------- CONFIG ----------------

const LISTEN_ADDR = "127.0.0.1:43595"

// Preload archive mapping (317)
// NOTE: title is NOT in idx0
var preloadArchives = []struct {
	Path    string
	IndexID int
	Archive int
}{
	{"/title", 1, 0}, // idx1, archive 0
	{"/config", 0, 1},
	{"/interface", 0, 2},
	{"/media", 0, 3},
	{"/versionlist", 0, 4},
	{"/textures", 0, 5},
	{"/wordenc", 0, 6},
	{"/sounds", 0, 7},
	{"/models", 0, 8},
}

// ---------------- CACHE ----------------

type Cache struct {
	data *os.File
	idx  map[int]*os.File
}

func openCache() (*Cache, error) {
	data, err := os.Open("cache/main_file_cache.dat")
	if err != nil {
		return nil, err
	}

	idx := make(map[int]*os.File)
	for i := 0; i <= 4; i++ {
		f, err := os.Open(fmt.Sprintf("cache/main_file_cache.idx%d", i))
		if err != nil {
			continue
		}
		idx[i] = f
	}

	return &Cache{data: data, idx: idx}, nil
}

// Extract reads a group from the cache using standard 317 sector chaining
func (c *Cache) Extract(indexID, archiveID int) ([]byte, error) {
	idx, ok := c.idx[indexID]
	if !ok {
		return nil, os.ErrNotExist
	}

	// idx entries are 6 bytes each
	_, err := idx.Seek(int64(archiveID*6), io.SeekStart)
	if err != nil {
		return nil, err
	}

	entry := make([]byte, 6)
	if _, err := io.ReadFull(idx, entry); err != nil {
		return nil, err
	}

	size := int(entry[0])<<16 | int(entry[1])<<8 | int(entry[2])
	sector := int(entry[3])<<16 | int(entry[4])<<8 | int(entry[5])
	if size == 0 || sector == 0 {
		return nil, os.ErrNotExist
	}

	var out bytes.Buffer
	read := 0
	chunk := 0

	for sector != 0 {
		log.Printf("idx%d[%d]: size=%d sector=%d", indexID, archiveID, size, sector)

		if _, err := c.data.Seek(int64(sector*520), io.SeekStart); err != nil {
			return nil, err
		}

		header := make([]byte, 8)
		if _, err := io.ReadFull(c.data, header); err != nil {
			return nil, err
		}

		curArchive := int(binary.BigEndian.Uint16(header[0:2]))
		curChunk := int(binary.BigEndian.Uint16(header[2:4]))
		nextSector := int(header[4])<<16 | int(header[5])<<8 | int(header[6])

		if curArchive != archiveID || curChunk != chunk {
			return nil, io.ErrUnexpectedEOF
		}

		remaining := size - read
		n := 512
		if remaining < n {
			n = remaining
		}

		payload := make([]byte, 512)
		if _, err := io.ReadFull(c.data, payload); err != nil {
			return nil, err
		}

		out.Write(payload[:n])
		read += n
		sector = nextSector
		chunk++
	}

	if read != size {
		return nil, io.ErrUnexpectedEOF
	}

	return out.Bytes(), nil
}

// ---------------- CRC TABLE ----------------

func build317CrcTable(cache *Cache) []byte {
	var crcs []uint32

	// CRC table excludes title (slot 0 is always zero)
	crcs = append(crcs, 0)

	for _, a := range preloadArchives[1:] {
		data, err := cache.Extract(a.IndexID, a.Archive)
		if err != nil {
			log.Printf("CRC extract failed for %s: %v", a.Path, err)
			crcs = append(crcs, 0)
			continue
		}
		crc := crc32.ChecksumIEEE(data)
		crcs = append(crcs, crc)
		log.Printf("Archive %d CRC = 0x%08X (%d bytes)", a.Archive, crc, len(data))
	}

	var buf bytes.Buffer
	for _, crc := range crcs {
		binary.Write(&buf, binary.BigEndian, crc)
	}

	hash := uint32(1234)
	for _, crc := range crcs {
		hash = (hash << 1) + crc
	}
	binary.Write(&buf, binary.BigEndian, hash)

	log.Printf("CRC table hash = 0x%08X", hash)
	return buf.Bytes()
}

func writeJaggrab(conn net.Conn, data []byte) error {
	var hdr [4]byte
	binary.BigEndian.PutUint32(hdr[:], uint32(len(data)))
	if _, err := conn.Write(hdr[:]); err != nil {
		return err
	}
	_, err := conn.Write(data)
	return err
}

// ---------------- SERVER ----------------

func handleConn(cache *Cache, conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	line = strings.TrimSpace(line)
	log.Println("JAGGRAB:", line)

	if !strings.HasPrefix(line, "JAGGRAB ") {
		return
	}

	path := strings.TrimPrefix(line, "JAGGRAB ")
	if i := strings.IndexByte(path, '-'); i != -1 {
		path = path[:i]
	}

	// Consume blank line
	reader.ReadString('\n')

	if strings.HasPrefix(path, "/crc") {
		data := build317CrcTable(cache)
		conn.Write(data)
		return
	}

	for _, a := range preloadArchives {
		if path == "/title" {
			data, err := cache.Extract(1, 0)
			if err != nil {
				log.Println("Title extract failed:", err)
				return
			}

			// Strip idx1 2-byte prefix
			if len(data) > 2 {
				data = data[2:]
			}

			log.Printf("Serving title (%d bytes)", len(data))
			writeJaggrab(conn, data)
			return
		}

		if path == a.Path {
			data, err := cache.Extract(a.IndexID, a.Archive)
			if err != nil {
				log.Println("Extract failed:", err)
				return
			}
			log.Printf("Serving %s (%d bytes)", path, len(data))
			conn.Write(data)
			return
		}
	}

	log.Println("Unknown JAGGRAB path:", path)
}

// ---------------- MAIN ----------------

func main() {
	cache, err := openCache()
	if err != nil {
		log.Fatal("Failed to open cache:", err)
	}
	defer cache.data.Close()
	for _, f := range cache.idx {
		defer f.Close()
	}

	ln, err := net.Listen("tcp", LISTEN_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Println("317 JAGGRAB server listening on", LISTEN_ADDR)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(cache, conn)
	}
}
