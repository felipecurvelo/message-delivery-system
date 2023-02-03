package main

import (
	"os"
	"sync"
	"time"

	"github.com/felipecurvelo/message-delivery-system/internal/server"
)

func main() {
	port := os.Args[1:]

	s := server.NewServer(&sync.Mutex{})
	s.Start(port[0])

	defer s.Close()

	for {
		time.Sleep(time.Millisecond * 300)
	}
}
