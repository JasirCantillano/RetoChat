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
		log.Fatalf("No se pudo iniciar el servidor: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("Servidor iniciado en puerto 8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("No se pudo conectar: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}

//Terminado
