package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/felipecurvelo/message-delivery-system/internal/client"
)

func main() {
	svr := "127.0.0.1:1234"

	c := client.NewClient()
	err := c.Connect(svr)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	clientOutputChan := make(chan string)
	go c.HandleMessages(clientOutputChan)

	c.WhoAmI()

	go func() {
		for {
			r := bufio.NewReader(os.Stdin)
			txt, _ := r.ReadString('\n')

			if txt == "list\n" {
				err := c.List()
				if err != nil {
					panic(err)
				}
			}

			if txt == "whoami\n" {
				err := c.WhoAmI()
				if err != nil {
					panic(err)
				}
			}

			spl := strings.Split(txt, "|")

			if len(spl) > 1 {
				c.Send(spl[0], []byte(spl[1]))
			}

			fmt.Printf("%v@%s> ", c.Id, svr)

			time.Sleep(time.Millisecond * 300)
		}
	}()

	for {
		if c.Id == "" {
			continue
		}

		fmt.Printf("%v@%s> ", c.Id, svr)
		msg := <-clientOutputChan
		fmt.Println(msg)

		time.Sleep(time.Millisecond * 300)
	}
}
