package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	name     string
	channels map[string]*channels
	commands chan<- commands
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
		case "/name":
			c.commands <- commands{
				id:     CMD_NAME,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- commands{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/channels":
			c.commands <- commands{
				id:     CMD_CHANNELS,
				client: c,
				args:   args,
			}
		case "/send":
			c.commands <- commands{
				id:     CMD_SEND,
				client: c,
				args:   args,
			}
		case "/exit":
			c.commands <- commands{
				id:     CMD_EXIT,
				client: c,
				args:   args,
			}
		case "/exitChannel":
			c.commands <- commands{
				id:     CMD_EXITCHANNEL,
				client: c,
				args:   args,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s ", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
