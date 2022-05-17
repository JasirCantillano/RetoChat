package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("Server started on port 8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Could not connect: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
