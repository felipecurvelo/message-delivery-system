package client

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Client struct {
	conn net.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(server string) error {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println(err)
		return err
	}
	c.conn = conn

	return nil
}

func (c *Client) SendMessage(clientID string, msg string) error {
	_, err := c.conn.Write([]byte(fmt.Sprintf("%s|%s", clientID, msg)))
	return err
}

func (c *Client) HandleMessages(clientOutputChan chan<- string) {
	for {
		c.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
		message, _ := bufio.NewReader(c.conn).ReadString('\n')
		if message != "" {
			clientOutputChan <- message
		}
	}
}
