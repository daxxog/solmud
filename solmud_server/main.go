// main.go
package main

import (
	"fmt"
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
