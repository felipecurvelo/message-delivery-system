package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	listener    net.Listener
	connections map[int]net.Conn
}

const maxClientID = 99

func NewServer() *Server {
	return &Server{
		connections: map[int]net.Conn{},
	}
}

func (s *Server) Start(port string) {
	s.connect(port)

	fmt.Println("Server started...")

	go s.WaitForNewConnections()

	go func() {
		for {
			s.HandleMessages()
		}
	}()
}

func (s *Server) connect(port string) {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.listener = l
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) WaitForNewConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		clientID := 1
		for {
			if _, ok := s.connections[clientID]; !ok {
				break
			}
			clientID++
		}

		s.connections[clientID] = conn
		fmt.Println("Client connected! Total clients connected: ", len(s.connections))
	}
}

func (s *Server) HandleMessages() {
	for connID := range s.connections {
		s.connections[connID].SetReadDeadline(time.Now().Add(time.Millisecond * 200))
		request, err := bufio.NewReader(s.connections[connID]).ReadString('\n')
		if err != nil {
			if opError, isOpError := err.(*net.OpError); isOpError && opError.Timeout() {
				continue
			}
			fmt.Println(err)
			return
		}

		r := strings.Split(request, "|")
		clientID, err := strconv.Atoi(r[0])
		if err != nil {
			fmt.Print("Invalid client id")
		}

		outputMsg := fmt.Sprintf("from:%v to:%v msg:%s", connID, clientID, string(r[1]))

		fmt.Print(">>> ", outputMsg)

		s.connections[clientID].Write([]byte(outputMsg))
	}
}
