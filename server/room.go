package server

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type Room struct {
	name    string
	clients map[net.Addr]*Client
}

func NewRoom(name string) *Room {
	return &Room{
		name:    name,
		clients: make(map[net.Addr]*Client),
	}
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) AddClient(client *Client) {
	if _, ok := r.clients[client.addr]; ok {
		log.Printf("client %s already exists", client.addr)
		return
	}
	r.clients[client.addr] = client

	// notify peers
	err := r.Welcome(client)
	if err != nil {
		log.Println(err)
		return
	}
}

func (r *Room) RemoveClient(client *Client) {
	if _, ok := r.clients[client.addr]; ok {
		delete(r.clients, client.addr)
	}
}

func (r *Room) ClientExist(client *Client) bool {
	if _, ok := r.clients[client.addr]; ok {
		return true
	}
	return false
}

func (r *Room) Welcome(client *Client) error {
	if client == nil {
		return errors.New("client is nil")
	}

	msg := fmt.Sprintf("> Welcome %s!\n", client.name)
	return r.Broadcast(msg, client)
}

func (r *Room) Broadcast(msg string, except *Client) error {
	for _, client := range r.clients {
		if except != nil && client.addr == except.addr {
			continue
		}
		_, err := client.conn.Write([]byte(msg))
		return err
	}
	return nil
}
