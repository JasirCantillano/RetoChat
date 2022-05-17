package main

type commandsID int

const (
	CMD_NAME commandsID = iota
	CMD_JOIN
	CMD_CHANNELS
	CMD_SEND
	CMD_EXIT
	CMD_EXITCHANNEL
)

type commands struct {
	id     commandsID
	client *client
	args   []string
}
