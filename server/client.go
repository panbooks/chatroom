package server

import (
	"net"
)

type Client struct {
	name string
	addr net.Addr
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		name: "anonymous",
		addr: conn.RemoteAddr(),
		conn: conn,
	}
}

func (c *Client) SetNick(name string) {
	c.name = name
}
