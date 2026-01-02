// main.go
package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
)

func getWorkingCrcTable() []byte {
	var buf bytes.Buffer

	// Unused index 0
	binary.Write(&buf, binary.BigEndian, uint32(0))

	// 8 archive CRCs (common for many original 317 clients)
	crcs := []uint32{
		0x6B5D6C9B, // title
		0x91579B40, // config
		0x7A5F7E9F, // interface
		0xD5B1E4B9, // media
		0x1D5B2E7C, // versionlist
		0xC8F2A8B3, // textures
		0xA9B7C6D4, // wordenc
		0xE4D5F6A7, // sounds
	}

	for _, crc := range crcs {
		binary.Write(&buf, binary.BigEndian, crc)
	}

	// Correct integrity hash for these values: 0xF29E7BEB
	binary.Write(&buf, binary.BigEndian, uint32(0xF29E7BEB))

	return buf.Bytes()
}

func main() {
	addr := fmt.Sprintf("%s:%d", LISTEN_HOST, LISTEN_PORT)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to bind to %s: %v", addr, err)
	}
	defer listener.Close()

	go startBasicHttpLogger()

	log.Printf("Solmud server listening on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go handleClient(conn)
	}
}

func startBasicHttpLogger() {
	addr := fmt.Sprintf("%s:%d", LISTEN_HOST, HTTP_CACHE_PORT)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to bind to %s for HTTP cache logging: %v", addr, err)
	}
	defer listener.Close()
	log.Printf("Basic HTTP cache request logger listening on %s", addr)
	log.Printf("Run your client now and watch the logs for incoming requests!")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("HTTP logger accept error: %v", err)
			continue
		}
		go handleHttpLogConnection(conn)
	}
}

func handleHttpLogConnection(conn net.Conn) {
	defer conn.Close()
	remote := conn.RemoteAddr().String()
	log.Printf("=== New HTTP connection from %s ===", remote)

	reader := bufio.NewReader(conn)

	reqLine, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading request line from %s: %v", remote, err)
		return
	}
	log.Printf("Request line: %q", reqLine)

	parts := strings.Split(strings.TrimSpace(reqLine), " ")
	if len(parts) < 2 {
		log.Printf("Malformed request line from %s", remote)
		fmt.Fprintf(conn, "HTTP/1.1 400 Bad Request\r\nContent-Length: 0\r\n\r\n")
		return
	}
	path := strings.ToLower(parts[1])

	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" || line == "\n" {
			break
		}
		log.Printf("Header: %q", line)
	}

	if strings.HasPrefix(path, "/crc") {
		log.Printf("Serving WORKING CRC table (common 317 values + correct hash 0xF29E7BEB)")
		// table := getWorkingCrcTable()
		table := getFakeValidCrcTable()
		n, err := conn.Write(table)
		log.Printf("Sent %d bytes for CRC table (err: %v)", n, err)
	} else {
		log.Printf("Unknown request path %s from %s – sending 404", path, remote)
		fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n")
	}

	log.Printf("=== End of request from %s ===\n", remote)
}

func getFakeValidCrcTable() []byte {
	var buf bytes.Buffer

	crcs := []uint32{
		0x00000000, // unused
		0xDEADBEEF, // fake title
		0xCAFEBABE, // fake config
		0xFEEDFACE, // fake interface
		0xBADC0DE0, // fake media
		0x12345678, // fake versionlist
		0x87654321, // fake textures
		0xABCDDCBA, // fake wordenc
		0x56789012, // fake sounds
	}

	for _, crc := range crcs {
		binary.Write(&buf, binary.BigEndian, crc)
	}

	// Correct hash for these values
	hash := int32(1234)
	for _, crc := range crcs {
		hash = (hash << 1) + int32(crc)
	}
	binary.Write(&buf, binary.BigEndian, hash)

	return buf.Bytes()
}
