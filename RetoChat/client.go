package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nombre   string
	canales  map[string]*canales
	comandos chan<- comandos
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nombre":
			c.comandos <- comandos{
				id:     CMD_NOMBRE,
				client: c,
				args:   args,
			}
		case "/subscribe":
			c.comandos <- comandos{
				id:     CMD_SUBSCRIBE,
				client: c,
				args:   args,
			}
		case "/canales":
			c.comandos <- comandos{
				id:     CMD_CANALES,
				client: c,
				args:   args,
			}
		case "/enviar":
			c.comandos <- comandos{
				id:     CMD_SEND,
				client: c,
				args:   args,
			}
		case "/salir":
			c.comandos <- comandos{
				id:     CMD_SALIR,
				client: c,
				args:   args,
			}
		case "/salirCanal":
			c.comandos <- comandos{
				id:     CMD_SALIRCANAL,
				client: c,
				args:   args,
			}
		default:
			c.err(fmt.Errorf("comando desconocido: %s ", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}

//Terminado
