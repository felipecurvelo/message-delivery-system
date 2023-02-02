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

func (c *Client) SendMessage(destinationIDs string, msg []byte) error {
	if len(destinationIDs) == 0 {
		return errors.New("Destination is empty")
	}

	if len(msg) == 0 {
		return errors.New("Message is empty")
	}

	if len(destinationIDs) > 255 {
		return errors.New("Max destination per message reached")
	}

	if len(msg) > 1024000 {
		return errors.New("Message size should be less than 1024Kb")
	}

	return c.sendMessageWithType("relay", destinationIDs, msg)
}

func (c *Client) WhoAmI() error {
	return c.sendMessageWithType("identity", "", []byte(""))
}

func (c *Client) List() error {
	return c.sendMessageWithType("list", "", []byte(""))
}

func (c *Client) sendMessageWithType(messageType string, destinationID string, msg []byte) error {
	inputMsg := append(msg, []byte("\n")...)
	_, err := c.conn.Write(append([]byte(fmt.Sprintf("%s|%s|", messageType, destinationID)), inputMsg...))
	return err
}

func (c *Client) HandleMessages(clientOutputChan chan<- string) error {
	for {
		c.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
		message, _ := c.reader.ReadString('\n')

		if strings.HasPrefix(message, "whoami:") {
			s := strings.Split(message, ":")
			c.Id = s[1]
			continue
		}

		if message != "" {
			clientOutputChan <- strings.Trim(message, "\n")
		}
	}
}
