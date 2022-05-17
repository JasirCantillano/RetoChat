package main

import "net"

type channels struct {
	name    string
	members map[net.Addr]*client
	sending map[int]*files
}

func (c *channels) broadcast(sender *client, msg string) {
	for addr, m := range c.members {
		if addr != sender.conn.RemoteAddr() {
			m.msg(msg)
		}
	}
}
