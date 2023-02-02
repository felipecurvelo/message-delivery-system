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
	address := "127.0.0.1:1234"

	client := client.NewClient()
	err := client.Connect(address)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	clientOutChan := make(chan string)
	go client.HandleMessages(clientOutChan)

	client.WhoAmI()

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			texts := strings.Split(text, "|")
			client.SendMessage(texts[0], []byte(texts[1]))

			fmt.Printf("%v@%s> ", client.Id, address)

			time.Sleep(time.Millisecond * 300)
		}
	}()

	for {
		if client.Id == "" {
			continue
		}

		fmt.Printf("%v@%s> ", client.Id, address)
		msg := <-clientOutChan
		fmt.Println(msg)

		time.Sleep(time.Millisecond * 300)
	}
}
