package client

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
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
	c.reader = bufio.NewReader(c.conn)

	return nil
}

func (s *Client) Close() {
	s.conn.Close()
}

func (c *Client) SendMessage(destinationID string, msg string) error {
	_, err := c.conn.Write([]byte(fmt.Sprintf("%s|%s", destinationID, msg)))
	fmt.Print("Message sent to ", destinationID)
	return err
}

func (c *Client) HandleMessages(clientOutputChan chan<- string) {
	for {
		c.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
		message, _ := c.reader.ReadString('\n')
		if message != "" {
			clientOutputChan <- message
		}
	}
}
