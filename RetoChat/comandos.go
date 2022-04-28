package main

type comandosID int

const (
	CMD_NOMBRE comandosID = iota
	CMD_SUBSCRIBE
	CMD_CANALES
	CMD_SEND
	CMD_SALIR
)

type comandos struct {
	id     comandosID
	client *client
	args   []string
}
