package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	m           *sync.Mutex
	listener    net.Listener
	connections map[int]net.Conn
}

const maxClientID = 99
const waitTimeDuration = time.Millisecond * 200

func NewServer(m *sync.Mutex) *Server {
	return &Server{
		connections: map[int]net.Conn{},
		m:           m,
	}
}

func (s *Server) Start(port string) error {
	err := s.connect(port)
	if err != nil {
		return err
	}

	fmt.Println("Server started...")

	go s.WaitForNewConnections()

	go func() {
		for {
			s.HandleMessages()
			s.Wait(waitTimeDuration)
		}
	}()
	return nil
}

func (s *Server) connect(port string) error {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return err
	}
	s.listener = l
	return nil
}

func (s *Server) Wait(duration time.Duration) {
	time.Sleep(duration)
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

		s.m.Lock()
		s.connections[clientID] = conn
		s.m.Unlock()
		fmt.Println("Client connected! Total clients connected: ", len(s.connections))
	}
}

func (s *Server) HandleMessages() {
	s.m.Lock()
	defer s.m.Unlock()
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
			outputMsg := fmt.Sprintf("whoami:%s", strconv.Itoa(sourceID))
			_, err = s.connections[sourceID].Write([]byte(outputMsg))
			break
		case "list":
			keys := []string{}
			for k := range s.connections {
				if k != sourceID {
					keys = append(keys, strconv.Itoa(k))
				}
			}
			sort.Strings(keys)
			outputMsg := strings.Join(keys, ",")
			_, err = s.connections[sourceID].Write([]byte(outputMsg))
			break
		case "relay":
			destinationIDs := strings.Split(r[1], ",")
			for _, id := range destinationIDs {
				destinationID, err := strconv.Atoi(id)
				if err != nil {
					fmt.Print("Invalid client id")
				}

				outputMsg := fmt.Sprintf("from:%v to:%v msg:%s", sourceID, destinationID, string(r[2]))

				_, err = s.connections[destinationID].Write([]byte(outputMsg))
				if err != nil {
					fmt.Println("ERROR", err)
				}
				fmt.Println(">>> ", outputMsg)
			}
			break
		}
	}
}
