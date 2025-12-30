// main.go
package main

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	addr := fmt.Sprintf("%s:%d", LISTEN_HOST, LISTEN_PORT)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to bind to %s: %v", addr, err)
	}
	defer listener.Close()

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

func handleClient(conn net.Conn) {
	defer conn.Close()
	log.Printf("New connection from %s", conn.RemoteAddr())

	// 1. Read handshake: opcode 14 + name hash
	header := make([]byte, 2)
	if _, err := io.ReadFull(conn, header); err != nil {
		log.Printf("Handshake read failed: %v", err)
		return
	}
	if header[0] != OPCODE_HANDSHAKE {
		log.Printf("Invalid handshake opcode: %d", header[0])
		return
	}
	nameHash := header[1]
	log.Printf("Handshake OK - nameHash: %d", nameHash)

	// 2. Generate random server session key
	var serverKey int64
	if err := binary.Read(cryptoRand.Reader, binary.BigEndian, &serverKey); err != nil {
		log.Printf("Random key generation failed, using 0: %v", err)
		serverKey = 0
	}

	// 3. Send 9-byte handshake response: 8-byte key + status 0
	resp := make([]byte, 9)
	binary.BigEndian.PutUint64(resp[:8], uint64(serverKey))
	resp[8] = 0 // status = proceed

	if _, err := conn.Write(resp); err != nil {
		log.Printf("Failed to send 9-byte handshake response: %v", err)
		return
	}
	log.Printf("Sent 9-byte handshake response (key: %d, status: 0)", serverKey)

	// 4. Send the 8-byte key AGAIN for flushInputStream(inStream.buffer, 8)
	if _, err := conn.Write(resp[:8]); err != nil {
		log.Printf("Failed to send extra 8-byte key: %v", err)
		return
	}
	log.Printf("Sent extra 8-byte key for client flushInputStream")

	// 5. Now read the login block header (type + payload size)
	loginHeader := make([]byte, 2)
	if _, err := io.ReadFull(conn, loginHeader); err != nil {
		log.Printf("Failed to read login block header: %v", err)
		return
	}

	loginType := loginHeader[0] // 16 = new login, 18 = reconnect
	payloadSize := int(loginHeader[1])

	log.Printf("Login block header received — type: %d, payload size: %d bytes", loginType, payloadSize)

	// 6. Read the encrypted payload
	payload := make([]byte, payloadSize)
	if _, err := io.ReadFull(conn, payload); err != nil {
		log.Printf("Failed to read encrypted payload: %v", err)
		return
	}
	log.Printf("Full encrypted login block received (%d bytes payload)", payloadSize)

	// 7. Send login success response: code 2
	successResponse := []byte{LOGIN_RESPONSE_OK, 0, 0}
	if _, err := conn.Write(successResponse); err != nil {
		log.Printf("Failed to send login success response: %v", err)
		return
	}

	log.Printf(">>> LOGIN SUCCESSFUL! Player 'daxxog' has logged in <<<")
	log.Printf("Client should now display black/loading screen")
	log.Printf("Connection remains open — ready for ISAAC cipher and map region packet")
}
