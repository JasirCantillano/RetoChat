package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	canales  map[string]*canales
	comandos chan comandos
}

func newServer() *server {
	return &server{
		canales:  make(map[string]*canales),
		comandos: make(chan comandos),
	}
}

func (s *server) run() {
	for cmd := range s.comandos {
		switch cmd.id {
		case CMD_NOMBRE:
			s.nombre(cmd.client, cmd.args)
		case CMD_SUBSCRIBE:
			s.subscribe(cmd.client, cmd.args)
		case CMD_CANALES:
			s.listarCanales(cmd.client, cmd.args)
		case CMD_SEND:
			s.send(cmd.client, cmd.args)
		case CMD_SALIR:
			s.salir(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("Nuevo cliente conectado: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nombre:   "anonimo",
		comandos: s.comandos,
	}

	c.readInput()
}

func (s *server) nombre(c *client, args []string) {
	c.nombre = args[1]
	c.msg(fmt.Sprintf("Hola, hola %s", c.nombre))
}

func (s *server) subscribe(c *client, args []string) {
	nombreCanal := args[1]
	r, ok := s.canales[nombreCanal]
	if !ok {
		r = &canales{
			nombre:   nombreCanal,
			miembros: make(map[net.Addr]*client),
		}
		s.canales[nombreCanal] = r
	}

	r.miembros[c.conn.RemoteAddr()] = c

	s.salirCanal(c)

	c.canales = r

	r.broadcast(c, fmt.Sprintf("%s se ha unido al canal", c.nombre))
	c.msg(fmt.Sprintf("Bienvenido %s", r.nombre))
}

func (s *server) listarCanales(c *client, args []string) {
	var canales []string
	for nombre := range s.canales {
		canales = append(canales, nombre)
	}

	c.msg(fmt.Sprintf("Los canales disponibles son: %s", strings.Join(canales, ", ")))
}

func (s *server) send(c *client, args []string) {
	if c.canales == nil {
		c.err(errors.New("Primero debes unirte al canal"))
		return
	}

	c.canales.broadcast(c, c.nombre+": "+strings.Join(args[1:len(args)], " "))
}

func (s *server) salir(c *client, args []string) {
	log.Printf("Un cliente se desconecto: %s", c.conn.RemoteAddr().String())

	s.salirCanal(c)

	c.msg("Esperamos tu regreso :(")
	c.conn.Close()
}

func (s *server) salirCanal(c *client) {
	if c.canales != nil {
		delete(c.canales.miembros, c.conn.RemoteAddr())
		c.canales.broadcast(c, fmt.Sprintf("%s Salio del canal", c.nombre))
	}
}
