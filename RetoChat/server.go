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
		case CMD_SALIRCANAL:
			s.salirCanal(cmd.client, cmd.args)
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
		canales:  make(map[string]*canales),
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
		c.canales[nombreCanal] = r
	}

	r.miembros[c.conn.RemoteAddr()] = c

	s.canales[nombreCanal] = r
	c.canales[nombreCanal] = r

	r.broadcast(c, fmt.Sprintf("%[1]s se ha unido al canal %[2]s", c.nombre, args[1]))
	c.msg(fmt.Sprintf("Bienvenido %[1]s al canal %[2]s", c.nombre, args[1]))
}

func (s *server) listarCanales(c *client, args []string) {
	var canales []string
	for nombre := range s.canales {
		canales = append(canales, nombre)
	}

	c.msg(fmt.Sprintf("Los canales disponibles son: %s", strings.Join(canales, ", ")))
}

func (s *server) send(c *client, args []string) {
	nombreCanal := args[1]

	var contador1 int = 0
	for k := range c.canales {
		if c.canales[k].nombre == nombreCanal {
			c.canales[nombreCanal].broadcast(c, c.nombre+" al canal "+args[1]+": "+strings.Join(args[2:len(args)], " "))
			contador1++
		}
	}
	if contador1 == 0 {
		c.err(errors.New("primero debes unirte al canal"))
		return
	}
}

func (s *server) salir(c *client, args []string) {
	log.Printf("Un cliente se desconecto: %s", c.conn.RemoteAddr().String())

	for k := range c.canales {
		viejoCanal := c.canales[k]
		delete(s.canales[k].miembros, c.conn.RemoteAddr())
		delete(c.canales[k].miembros, c.conn.RemoteAddr())
		delete(c.canales, k)
		var mensaje string = "Esperamos tu regreso al canal " + args[1] + " muy pronto :("
		c.msg(mensaje)
		viejoCanal.broadcast(c, fmt.Sprintf("%[1]s Salio del canal %[2]s", c.nombre, args[1]))
	}

	c.msg("Esperamos tu regreso :(")
	c.conn.Close()
}

func (s *server) salirCanal(c *client, args []string) {
	nombreCanal := args[1]
	viejoCanal := c.canales[nombreCanal]
	for k := range c.canales {
		if c.canales[k].nombre == nombreCanal {
			delete(s.canales[nombreCanal].miembros, c.conn.RemoteAddr())
			delete(c.canales[nombreCanal].miembros, c.conn.RemoteAddr())
			delete(c.canales, nombreCanal)
			var mensaje string = "Esperamos tu regreso al canal " + args[1] + " muy pronto :("
			c.msg(mensaje)
			viejoCanal.broadcast(c, fmt.Sprintf("%[1]s Salio del canal %[2]s", c.nombre, args[1]))
		}
	}
}
