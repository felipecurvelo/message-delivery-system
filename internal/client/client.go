package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	Id     string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(server string) error {
	r := strings.Split(server, ":")
	if len(r) < 2 {
		return errors.New("Please provide server:port")
	}

	conn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println(err)
		return err
	}
	c.conn = conn
	c.reader = bufio.NewReader(c.conn)

	return nil
}

func (s *Client) Close() {
	s.conn.Close()
}

func (c *Client) Send(destination string, msg []byte) error {
	if len(destination) == 0 {
		return errors.New("Destination is empty")
	}

	if len(msg) == 0 {
		return errors.New("Message is empty")
	}

	if len(destination) > 255 {
		return errors.New("Max destination per message reached")
	}

	if len(msg) > 1024000 {
		return errors.New("Message size should be less than 1024Kb")
	}

	return c.sendWithType("send", destination, msg)
}

func (c *Client) WhoAmI() error {
	return c.sendWithType("whoami", "", []byte(""))
}

func (c *Client) List() error {
	return c.sendWithType("list", "", []byte(""))
}

func (c *Client) sendWithType(messageType string, destinationID string, msg []byte) error {
	input := append(msg, []byte("\n")...)
	_, err := c.conn.Write(append([]byte(fmt.Sprintf("%s|%s|", messageType, destinationID)), input...))
	return err
}

func (c *Client) HandleMessages(clientOutputChan chan<- string) error {
	for {
		c.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
		msg, _ := c.reader.ReadString('\n')

		if strings.HasPrefix(msg, "whoami:") {
			split := strings.Split(msg, ":")
			c.Id = split[1]
		}

		if msg != "" {
			clientOutputChan <- strings.Trim(msg, "\n")
		}
	}
}
