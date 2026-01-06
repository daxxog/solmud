// main.go
package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"strings"
)

const (
	LISTEN_ADDR = "127.0.0.1:43595" // MUST match client.j(43595)
)

// ---- CRC TABLE ----
// 9 CRCs + 1 hash = 40 bytes total
func getValid317CrcTable() []byte {
	var buf bytes.Buffer

	crcs := []uint32{
		0x6B5D6C9B, // 1 title
		0x91579B40, // 2 config
		0x7A5F7E9F, // 3 interface
		0xD5B1E4B9, // 4 media
		0x1D5B2E7C, // 5 versionlist
		0xC8F2A8B3, // 6 textures
		0xA9B7C6D4, // 7 wordenc
		0xE4D5F6A7, // 8 sounds
		0x1A0A1B9C, // 9 models (critical – was 0)
	}

	for _, crc := range crcs {
		_ = binary.Write(&buf, binary.BigEndian, crc)
	}

	// Exact client hash computation (unsigned 32-bit)
	hash := uint32(1234)
	for _, crc := range crcs {
		hash = (hash << 1) + crc
	}
	_ = binary.Write(&buf, binary.BigEndian, hash)

	return buf.Bytes() // 40 bytes
}

// ---- CONNECTION HANDLER ----
func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read "JAGGRAB /crcXXXX-317\n"
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Println("read error:", err)
		return
	}

	line = strings.TrimSpace(line)
	log.Println("JAGGRAB:", line)

	if !strings.HasPrefix(line, "JAGGRAB /crc") {
		log.Println("Ignoring non-CRC request")
		return
	}

	// Consume blank line
	_, _ = reader.ReadString('\n')

	// Send raw CRC bytes ONLY
	table := getValid317CrcTable()
	n, err := conn.Write(table)
	log.Printf("Sent %d CRC bytes (err=%v)", n, err)
}

// ---- MAIN ----
func main() {
	log.Println("Starting 317 JAGGRAB CRC server on", LISTEN_ADDR)

	ln, err := net.Listen("tcp", LISTEN_ADDR)
	if err != nil {
		log.Fatal("listen failed:", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go handleConn(conn)
	}
}
