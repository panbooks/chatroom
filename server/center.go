package server

import (
	"fmt"
	"log"
	"net"
)

type Center struct {
	rooms   map[string]*Room
	clients map[net.Addr]*Client
}

func NewCenter() *Center {
	return &Center{
		rooms:   make(map[string]*Room),
		clients: make(map[net.Addr]*Client),
	}
}

func (c *Center) AddRoom(room *Room) {
	if _, ok := c.rooms[room.name]; ok {
		log.Printf("room %s already exists", room.name)
		return
	}
	c.rooms[room.name] = room
}

func (c *Center) RemoveRoom(room *Room) {
	if _, ok := c.rooms[room.name]; ok {
		delete(c.rooms, room.name)
	}
}

func (c *Center) GetRoom(roomName string) (*Room, error) {
	if _, ok := c.rooms[roomName]; !ok {
		log.Fatalf("room %s does not exists", roomName)
		return nil, fmt.Errorf("room %s does not exists", roomName)
	}
	return c.rooms[roomName], nil
}

func (c *Center) ListRooms() []string {
	roomNames := make([]string, 0, len(c.rooms))
	for _, room := range c.rooms {
		roomNames = append(roomNames, room.name)
	}
	return roomNames
}

func (c *Center) AddClient(client *Client) {
	if _, ok := c.clients[client.addr]; ok {
		log.Printf("client %s already exists", client.addr)
		return
	}
	c.clients[client.addr] = client
}

func (c *Center) RemoveClient(client *Client) {
	if _, ok := c.clients[client.addr]; ok {
		delete(c.clients, client.addr)
	}
}

func (c *Center) RoomExist(roomName string) bool {
	if _, ok := c.rooms[roomName]; ok {
		return true
	}
	return false
}

func (c *Center) ClientExist(addr net.Addr) bool {
	if _, ok := c.clients[addr]; ok {
		return true
	}
	return false
}
