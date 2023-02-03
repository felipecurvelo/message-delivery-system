package test

import (
	"sync"
	"testing"

	"github.com/felipecurvelo/message-delivery-system/internal/client"
	"github.com/felipecurvelo/message-delivery-system/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationIdentity(t *testing.T) {
	s := server.NewServer(&sync.Mutex{})
	s.Start("1234")
	defer s.Close()

	client1 := client.NewClient()
	err := client1.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client1.Close()

	client2 := client.NewClient()
	err = client2.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client2.Close()

	client1OutChan := make(chan string)
	go client1.HandleMessages(client1OutChan)

	client2OutChan := make(chan string)
	go client2.HandleMessages(client2OutChan)

	err = client1.WhoAmI()
	assert.NoError(t, err)

	err = client2.WhoAmI()
	assert.NoError(t, err)

	client1ID := <-client1OutChan
	assert.Equal(t, "whoami:1", client1ID)

	client2ID := <-client2OutChan
	assert.Equal(t, "whoami:2", client2ID)
}

func TestIntegrationList(t *testing.T) {
	s := server.NewServer(&sync.Mutex{})
	s.Start("1234")
	defer s.Close()

	client1 := client.NewClient()
	err := client1.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client1.Close()

	client2 := client.NewClient()
	err = client2.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client2.Close()

	client3 := client.NewClient()
	err = client3.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client3.Close()

	client1OutChan := make(chan string)
	go client1.HandleMessages(client1OutChan)

	client2OutChan := make(chan string)
	go client2.HandleMessages(client2OutChan)

	err = client1.List()
	assert.NoError(t, err)

	err = client2.List()
	assert.NoError(t, err)

	client1ID := <-client1OutChan
	assert.Equal(t, "list:[1,2,3]", client1ID)

	client2ID := <-client2OutChan
	assert.Equal(t, "list:[1,2,3]", client2ID)
}

func TestIntegrationSingleDestination(t *testing.T) {
	s := server.NewServer(&sync.Mutex{})
	s.Start("1234")
	defer s.Close()

	client1 := client.NewClient()
	err := client1.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client1.Close()

	client2 := client.NewClient()
	err = client2.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client2.Close()

	client1OutChan := make(chan string)
	go client1.HandleMessages(client1OutChan)

	client2OutChan := make(chan string)
	go client2.HandleMessages(client2OutChan)

	//Send a first msg
	err = client1.Send("2", []byte("firstMsgToClient2"))
	assert.NoError(t, err)

	err = client2.Send("1", []byte("firstMsgToClient1"))
	assert.NoError(t, err)

	firstMsgToClient2 := <-client2OutChan
	assert.Equal(t, "from:1 to:2 msg:firstMsgToClient2", firstMsgToClient2)

	firstMsgToClient1 := <-client1OutChan
	assert.Equal(t, "from:2 to:1 msg:firstMsgToClient1", firstMsgToClient1)

	// Send a second msg
	err = client1.Send("2", []byte("secondMsgToClient2"))
	assert.NoError(t, err)

	err = client2.Send("1", []byte("secondMsgToClient1"))
	assert.NoError(t, err)

	secondMsgToClient2 := <-client2OutChan
	assert.Equal(t, "from:1 to:2 msg:secondMsgToClient2", secondMsgToClient2)

	secondMsgToClient1 := <-client1OutChan
	assert.Equal(t, "from:2 to:1 msg:secondMsgToClient1", secondMsgToClient1)
}

func TestIntegrationMultipleDestination(t *testing.T) {
	s := server.NewServer(&sync.Mutex{})
	s.Start("1234")
	defer s.Close()

	client1 := client.NewClient()
	err := client1.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client1.Close()

	client2 := client.NewClient()
	err = client2.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client2.Close()

	client3 := client.NewClient()
	err = client3.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client3.Close()

	client1OutChan := make(chan string)
	go client1.HandleMessages(client1OutChan)

	client2OutChan := make(chan string)
	go client2.HandleMessages(client2OutChan)

	client3OutChan := make(chan string)
	go client3.HandleMessages(client3OutChan)

	//Send a first msg
	err = client1.Send("2,3", []byte("msg"))
	assert.NoError(t, err)

	msgToClient2 := <-client2OutChan
	assert.Equal(t, "from:1 to:2 msg:msg", msgToClient2)

	msgToClient3 := <-client3OutChan
	assert.Equal(t, "from:1 to:3 msg:msg", msgToClient3)
}

func TestIntegrationMultipleDestinationRaceCondition(t *testing.T) {
	s := server.NewServer(&sync.Mutex{})
	s.Start("1234")
	defer s.Close()

	client1 := client.NewClient()
	err := client1.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client1.Close()

	client2 := client.NewClient()
	err = client2.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client2.Close()

	client3 := client.NewClient()
	err = client3.Connect("127.0.0.1:1234")
	assert.NoError(t, err)
	defer client3.Close()

	client1OutChan := make(chan string)
	go client1.HandleMessages(client1OutChan)

	client2OutChan := make(chan string)
	go client2.HandleMessages(client2OutChan)

	client3OutChan := make(chan string)
	go client3.HandleMessages(client3OutChan)

	//Send a first msg
	err = client1.Send("2,3", []byte("msg"))
	assert.NoError(t, err)

	msgToClient2 := <-client2OutChan
	assert.Equal(t, "from:1 to:2 msg:msg", msgToClient2)

	msgToClient3 := <-client3OutChan
	assert.Equal(t, "from:1 to:3 msg:msg", msgToClient3)

	//Send a second msg
	err = client1.Send("2,3", []byte("msg"))
	assert.NoError(t, err)

	secondMsgToClient2 := <-client2OutChan
	assert.Equal(t, "from:1 to:2 msg:msg", secondMsgToClient2)

	secondMsgToClient3 := <-client3OutChan
	assert.Equal(t, "from:1 to:3 msg:msg", secondMsgToClient3)

	//Send a third msg
	err = client2.Send("1,3", []byte("msg"))
	assert.NoError(t, err)

	thirdMsgToClient1 := <-client1OutChan
	assert.Equal(t, "from:2 to:1 msg:msg", thirdMsgToClient1)

	thirdMsgToClient3 := <-client3OutChan
	assert.Equal(t, "from:2 to:3 msg:msg", thirdMsgToClient3)
}
