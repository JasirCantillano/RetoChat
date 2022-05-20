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
	var count int = 0
	for range args {
		count++
	}
	if count > 1 && len(args[1]) > 0 {
		c.name = args[1]
		c.msg(fmt.Sprintf("Hello, hello %s", c.name))
	} else {
		c.err(errors.New("a name was not defined"))
	}
}

func (s *server) join(c *client, args []string) {
	var count int = 0
	for range args {
		count++
	}
	if count > 1 && len(args[1]) > 0 {
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
		s.Statistic(w, v)
	} else {
		c.err(errors.New("a channel was not defined"))
	}

}

func (s *server) listChannels(c *client, args []string) {
	var channels []string
	var count int = 0
	for name := range s.channels {
		channels = append(channels, name)
		count++
	}
	if count > 0 {
		c.msg(fmt.Sprintf("Available channels are: %s", strings.Join(channels, ", ")))
	} else {
		c.err(errors.New("no channels yet"))
	}
}

func (s *server) send(c *client, args []string) {
	var count int = 0
	for range args {
		count++
	}
	if count > 2 {
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
				count1++
				var fileOrigin string = "files/" + args[2]
				origin, err := os.Open(fileOrigin)
				if err != nil {
					c.err(errors.New("the specified file does not exist"))
				} else {
					defer origin.Close()
					var fileDestination string = directory + "/" + args[2]
					destination, err := os.OpenFile(fileDestination, os.O_RDWR|os.O_CREATE, 0666)
					if err != nil {
						log.Fatal(err)
					}
					defer destination.Close()
					io.Copy(destination, origin)
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
		}
		if count1 == 0 {
			c.err(errors.New("you must first join a channel"))
			return
		}
	} else {
		c.err(errors.New("a channel or a file was not defined"))
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

	var w http.ResponseWriter
	var v *http.Request
	s.Start(w, v)
}

func (s *server) exitChannel(c *client, args []string) {
	var count int = 0
	for range args {
		count++
	}
	if count > 1 && len(args[1]) > 0 {
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

				var w http.ResponseWriter
				var v *http.Request
				s.Start(w, v)
			} else {
				c.err(errors.New("this channel does not exist"))
			}
		}
	} else {
		c.err(errors.New("a channel was not defined"))
	}
}
