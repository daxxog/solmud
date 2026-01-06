package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// ---------------- CONFIG ----------------

const LISTEN_ADDR = "127.0.0.1:43595"

// 317 preload archive IDs (idx0)
var preloadArchives = []struct {
	Path string
	ID   int
}{
	{"/title", 0},
	{"/config", 1},
	{"/interface", 2},
	{"/media", 3},
	{"/versionlist", 4},
	{"/textures", 5},
	{"/wordenc", 6},
	{"/sounds", 7},
	{"/models", 8},
}

// ---------------- CACHE ----------------

type Cache struct {
	data *os.File
	idx0 *os.File
}

func openCache() (*Cache, error) {
	data, err := os.Open("cache/main_file_cache.dat")
	if err != nil {
		return nil, err
	}
	idx0, err := os.Open("cache/main_file_cache.idx0")
	if err != nil {
		data.Close()
		return nil, err
	}
	return &Cache{data: data, idx0: idx0}, nil
}

// ExtractRawArchive reads an idx0 archive and returns the *exact compressed JAG bytes*
func (c *Cache) ExtractRawArchive(archiveID int) ([]byte, error) {
	// idx0 entry = 6 bytes
	_, err := c.idx0.Seek(int64(archiveID*6), io.SeekStart)
	if err != nil {
		return nil, err
	}

	idx := make([]byte, 6)
	if _, err := io.ReadFull(c.idx0, idx); err != nil {
		return nil, err
	}

	size := int(idx[0])<<16 | int(idx[1])<<8 | int(idx[2])
	sector := int(idx[3])<<16 | int(idx[4])<<8 | int(idx[5])
	if size == 0 || sector == 0 {
		return nil, os.ErrNotExist
	}

	var out bytes.Buffer
	read := 0
	chunk := 0

	for sector != 0 {
		log.Printf(
			"idx0[%d]: size=%d sector=%d",
			archiveID, size, sector,
		)

		// Each sector = 520 bytes
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

		// indexID := int(header[7])
		// if curArchive != archiveID || curChunk != chunk || indexID != 0 {
		if curArchive != archiveID || curChunk != chunk {
			return nil, io.ErrUnexpectedEOF
		}

		remaining := size - read
		dataSize := 512
		if remaining < dataSize {
			dataSize = remaining
		}

		payload := make([]byte, 512)
		if _, err := io.ReadFull(c.data, payload); err != nil {
			return nil, err
		}

		out.Write(payload[:dataSize])
		read += dataSize
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

	for _, a := range preloadArchives {
		data, err := cache.ExtractRawArchive(a.ID)
		if err != nil {
			log.Printf("CRC extract failed for %s: %v", a.Path, err)
			crcs = append(crcs, 0)
			continue
		}
		crc := crc32.ChecksumIEEE(data)
		crcs = append(crcs, crc)
		log.Printf("Archive %d CRC = 0x%08X (%d bytes)", a.ID, crc, len(data))
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

	// Strip optional -CRC suffix (e.g. /title-123456789)
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
		if path == a.Path {
			data, err := cache.ExtractRawArchive(a.ID)
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
	defer cache.idx0.Close()

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
