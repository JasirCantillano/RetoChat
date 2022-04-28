package main

import "net"

type canales struct {
	nombre   string
	miembros map[net.Addr]*client
}
