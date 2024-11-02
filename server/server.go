package server

import (
	"bytes"
	"log"
	"net"
)

type Server struct {
	center *Center
}

func NewServer() *Server {
	return &Server{
		center: NewCenter(),
	}
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on port 8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Accepted new connection from %s", conn.RemoteAddr())

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	client := NewClient(conn)

	buf := make([]byte, 64)
	for {
		log.Println("client name:", client.name)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		splits := bytes.Split(bytes.Trim(buf[:n], "\r\n"), []byte(" "))
		commandName := string(splits[0])
		arg := ""
		if len(splits) > 1 {
			arg = string(splits[1])
		}
		log.Printf("get command:%s, arg:%s", commandName, arg)
		command := NewCommand(commandName, arg, client, s.center)
		command.Execute()
	}
}
