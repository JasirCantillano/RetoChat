package main

import "net"

type canales struct {
	nombre   string
	miembros map[net.Addr]*client
	envios   map[int]*archivos
}

func (c *canales) broadcast(sender *client, msg string) {
	for addr, m := range c.miembros {
		if addr != sender.conn.RemoteAddr() {
			m.msg(msg)
		}
	}
}
