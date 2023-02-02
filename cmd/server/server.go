package main

import (
	"sync"
	"time"

	"github.com/felipecurvelo/message-delivery-system/internal/server"
)

func main() {
	port := "1234"

	s := server.NewServer(&sync.Mutex{})
	s.Start(port)

	defer s.Close()

	for {
		time.Sleep(time.Millisecond * 300)
	}
}
