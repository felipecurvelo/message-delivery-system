package server

import (
	"bufio"
	"fmt"
	"io"
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
	for sourceID := range s.connections {
		s.connections[sourceID].SetReadDeadline(time.Now().Add(time.Millisecond * 200))
		request, err := bufio.NewReader(s.connections[sourceID]).ReadString('\n')
		if err != nil {
			if opError, isOpError := err.(*net.OpError); isOpError && opError.Timeout() {
				continue
			}
			if err == io.EOF {
				delete(s.connections, sourceID)
				continue
			}
			fmt.Print("ERROR", err)
			return
		}

		r := strings.Split(request, "|")

		switch r[0] {
		case "identity":
			outputMsg := strconv.Itoa(sourceID)
			_, err = s.connections[sourceID].Write([]byte(outputMsg))
			break
		case "list":

			break
		case "relay":
			destinationIDs := strings.Split(r[1], ",")
			for _, id := range destinationIDs {
				destinationID, err := strconv.Atoi(id)
				if err != nil {
					fmt.Print("Invalid client id")
				}

				outputMsg := fmt.Sprintf("from:%v to:%v msg:%s", sourceID, destinationID, string(r[2]))

				fmt.Println("WRITING TO ", destinationID)

				_, err = s.connections[destinationID].Write([]byte(outputMsg))
				if err != nil {
					fmt.Print("ERROR", err)
				}
				fmt.Println(">>> ", outputMsg)
				fmt.Println("WROTE TO ", destinationID)
			}
			break
		}
	}
}
