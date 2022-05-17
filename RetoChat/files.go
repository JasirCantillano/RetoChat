package main

import "net"

type files struct {
	name      string
	client    string
	direction net.Addr
}
