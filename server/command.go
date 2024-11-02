package server

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type CommandName string

// NICK nick {nickname}
const NICK CommandName = "nick"

// JOIN join {roomname}
const JOIN CommandName = "join"

// LEAVE : leave room
const LEAVE CommandName = "leave"

// SAY {words}
const SAY CommandName = "say"

const LIST CommandName = "list"

type Command struct {
	name   CommandName
	arg    string
	client *Client
	center *Center
}

func NewCommand(name, arg string, client *Client, center *Center) *Command {
	return &Command{
		name:   CommandName(name),
		arg:    arg,
		client: client,
		center: center,
	}
}

func (c *Command) Execute() {
	var err error
	var res string
	switch c.name {
	case NICK:
		res, err = c.nick()
	case JOIN:
		res, err = c.join()
	case LEAVE:
		res, err = c.leave()
	case SAY:
		res, err = c.say()
	case LIST:
		res, err = c.list()
	default:
		log.Printf("Unknown command: %s", c.name)
		err = errors.New("unknown command")
	}
	if err != nil {
		c.client.conn.Write([]byte(fmt.Sprintf("> %s\n", err)))
	} else {
		c.client.conn.Write([]byte(fmt.Sprintf("> %s\n", res)))
	}
}

func (c *Command) nick() (string, error) {
	fmt.Println(c.client.name)
	if c.arg == "" {
		return c.client.name, nil
	}
	c.client.SetNick(c.arg)
	fmt.Println(c.client.name)
	return "ok", nil
}

func (c *Command) join() (string, error) {
	if c.arg == "" {
		return "", errors.New("param error for join")
	}

	roomName := c.arg
	c.center.AddRoom(NewRoom(roomName))

	room, err := c.center.GetRoom(roomName)
	if err != nil {
		return "", err
	}
	room.AddClient(c.client)
	c.center.AddClient(c.client)

	return "ok", nil
}

func (c *Command) leave() (string, error) {
	if c.arg == "" {
		return "", errors.New("need param {room} for leave")
	}

	roomName := c.arg
	if !c.center.RoomExist(roomName) {
		log.Printf("room %s does not exist", roomName)
		return "", errors.New(fmt.Sprintf("room %s does not exist", roomName))
	}

	// remove client in the room
	room, err := c.center.GetRoom(roomName)
	if err != nil {
		return "", err
	}
	room.RemoveClient(c.client)

	// notify to peers
	msg := fmt.Sprintf("> %s left room\n", c.client.name)
	room.Broadcast(msg, nil)

	// if room has no clients, delete the room
	if len(room.clients) == 0 {
		c.center.RemoveRoom(room)
	}

	return "ok", nil
}

func (c *Command) say() (string, error) {
	// one client can join several rooms
	// need to keep which room he is in
	if c.arg == "" {
		return "", errors.New("param error for say")
	}

	words := fmt.Sprintf("> %s: %s\n", c.client.name, c.arg)
	for _, room := range c.center.rooms {
		if room.ClientExist(c.client) {
			err := room.Broadcast(words, c.client)
			if err != nil {
				return "", err
			}
		}
	}

	return "ok", nil
}

func (c *Command) list() (string, error) {
	rooms := c.center.ListRooms()
	if len(rooms) == 0 {
		return "no rooms found", nil
	} else {
		return strings.Join(rooms, ", "), nil
	}
}

func (c *Command) String() string {
	return fmt.Sprintf("%s %s", c.name, c.arg)
}
