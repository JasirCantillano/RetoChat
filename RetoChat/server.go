package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type server struct {
	channels map[string]*channels
	commands chan commands
}

func newServer() *server {
	return &server{
		channels: make(map[string]*channels),
		commands: make(chan commands),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NAME:
			s.name(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_CHANNELS:
			s.listChannels(cmd.client, cmd.args)
		case CMD_SEND:
			s.send(cmd.client, cmd.args)
		case CMD_EXITCHANNEL:
			s.exitChannel(cmd.client, cmd.args)
		case CMD_EXIT:
			s.exit(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("New client conected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		name:     "anonymous",
		channels: make(map[string]*channels),
		commands: s.commands,
	}

	c.readInput()
}

func (s *server) name(c *client, args []string) {
	c.name = args[1]
	c.msg(fmt.Sprintf("Hello, hello %s", c.name))
}

func (s *server) join(c *client, args []string) {
	nameChannel := args[1]
	r, ok := s.channels[nameChannel]
	if !ok {
		r = &channels{
			name:    nameChannel,
			members: make(map[net.Addr]*client),
			sending: make(map[int]*files),
		}
		s.channels[nameChannel] = r
		c.channels[nameChannel] = r
	}

	r.members[c.conn.RemoteAddr()] = c

	s.channels[nameChannel] = r
	c.channels[nameChannel] = r

	r.broadcast(c, fmt.Sprintf("%[1]s has joined the channel %[2]s", c.name, args[1]))
	c.msg(fmt.Sprintf("Welcome %[1]s the channel %[2]s", c.name, args[1]))

	var w http.ResponseWriter
	var v *http.Request
	s.Start(w, v)
}

func (s *server) listChannels(c *client, args []string) {
	var channels []string
	for name := range s.channels {
		channels = append(channels, name)
	}

	c.msg(fmt.Sprintf("Available channels are: %s", strings.Join(channels, ", ")))
}

func (s *server) send(c *client, args []string) {
	nameChannel := args[1]

	var count1 int = 0
	var count2 int = 0
	for k := range c.channels {
		if c.channels[k].name == nameChannel {
			var directory string = "files/" + c.channels[k].name
			err := os.Mkdir(directory, 0750)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}
			var fileOrigin string = "files/" + args[2]
			origin, err := os.Open(fileOrigin)
			if err != nil {
				log.Fatal(err)
			}
			defer origin.Close()
			var fileDestination string = directory + "/" + args[2]
			destination, err := os.OpenFile(fileDestination, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Fatal(err)
			}
			defer destination.Close()
			io.Copy(destination, origin)
			count1++
			for range s.channels[nameChannel].sending {
				count2++
			}
			r, ok := s.channels[nameChannel].sending[count2]
			if !ok {
				r = &files{
					name:      args[2],
					client:    c.name,
					direction: c.conn.RemoteAddr(),
				}
				s.channels[nameChannel].sending[count2] = r
				c.channels[nameChannel].sending[count2] = r
			}
			c.channels[nameChannel].sending[count2] = r
			c.channels[nameChannel].broadcast(c, c.name+" the channel "+args[1]+": "+strings.Join(args[2:len(args)], " ")+" send to this channel, check it in the folder: "+directory)

			var w http.ResponseWriter
			var v *http.Request
			s.Statistic(w, v)
		}
	}
	if count1 == 0 {
		c.err(errors.New("you must first join a channel"))
		return
	}
}

func (s *server) exit(c *client, args []string) {
	log.Printf("A client disconnected: %s", c.conn.RemoteAddr().String())

	for k := range c.channels {
		oldChannel := c.channels[k]
		delete(s.channels[k].members, c.conn.RemoteAddr())
		delete(c.channels[k].members, c.conn.RemoteAddr())
		delete(c.channels, k)
		var messaje string = "We look forward to your return to the server" + k + " very soon :("
		c.msg(messaje)
		oldChannel.broadcast(c, fmt.Sprintf("%[1]s Left the channel %[2]s", c.name, k))
	}

	c.msg("We look forward to your return :(")
	c.conn.Close()
}

func (s *server) exitChannel(c *client, args []string) {
	nameChannel := args[1]
	oldChannel := c.channels[nameChannel]
	for k := range c.channels {
		if c.channels[k].name == nameChannel {
			delete(s.channels[nameChannel].members, c.conn.RemoteAddr())
			delete(c.channels[nameChannel].members, c.conn.RemoteAddr())
			delete(c.channels, nameChannel)
			var messaje string = "We look forward to your return to the channel " + args[1] + " very soon :("
			c.msg(messaje)
			oldChannel.broadcast(c, fmt.Sprintf("%[1]s Left the channel %[2]s", c.name, args[1]))
		}
	}
}
