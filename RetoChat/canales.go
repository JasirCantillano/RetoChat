package main

import "net"

type canales struct {
	nombre   string
	miembros map[net.Addr]*client
}

func (c *canales) broadcast(sender *client, msg string) {
	for addr, m := range c.miembros {
		if addr != sender.conn.RemoteAddr() {
			m.msg(msg)
		}
	}
}
