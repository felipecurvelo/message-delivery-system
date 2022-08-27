package test

import (
	"testing"
	"time"

	"github.com/felipecurvelo/message-delivery-system/internal/client"
	"github.com/felipecurvelo/message-delivery-system/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationIdentity(t *testing.T) {
	s := server.NewServer()
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

	time.Sleep(500 * time.Millisecond)

	client1OutChan := make(chan string)
	go client1.HandleMessages(client1OutChan)

	client2OutChan := make(chan string)
	go client2.HandleMessages(client2OutChan)

	err = client1.WhoAmI()
	assert.NoError(t, err)

	err = client2.WhoAmI()
	assert.NoError(t, err)

	client1ID := <-client1OutChan
	assert.Equal(t, "1", client1ID)

	client2ID := <-client2OutChan
	assert.Equal(t, "2", client2ID)
}

func TestIntegrationList(t *testing.T) {
	s := server.NewServer()
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

	time.Sleep(500 * time.Millisecond)

	client1OutChan := make(chan string)
	go client1.HandleMessages(client1OutChan)

	client2OutChan := make(chan string)
	go client2.HandleMessages(client2OutChan)

	err = client1.List()
	assert.NoError(t, err)

	err = client2.List()
	assert.NoError(t, err)

	client1ID := <-client1OutChan
	assert.Equal(t, "2,3", client1ID)

	client2ID := <-client2OutChan
	assert.Equal(t, "1,3", client2ID)
}

func TestIntegrationSingleDestination(t *testing.T) {
	s := server.NewServer()
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
	err = client1.SendMessage("2", "firstMsgToClient2\n")
	assert.NoError(t, err)

	err = client2.SendMessage("1", "firstMsgToClient1\n")
	assert.NoError(t, err)

	firstMsgToClient2 := <-client2OutChan
	assert.Equal(t, "from:1 to:2 msg:firstMsgToClient2\n", firstMsgToClient2)

	firstMsgToClient1 := <-client1OutChan
	assert.Equal(t, "from:2 to:1 msg:firstMsgToClient1\n", firstMsgToClient1)

	// Send a second msg
	err = client1.SendMessage("2", "secondMsgToClient2\n")
	assert.NoError(t, err)

	err = client2.SendMessage("1", "secondMsgToClient1\n")
	assert.NoError(t, err)

	secondMsgToClient2 := <-client2OutChan
	assert.Equal(t, "from:1 to:2 msg:secondMsgToClient2\n", secondMsgToClient2)

	secondMsgToClient1 := <-client1OutChan
	assert.Equal(t, "from:2 to:1 msg:secondMsgToClient1\n", secondMsgToClient1)
}

func TestIntegrationMultipleDestination(t *testing.T) {
	s := server.NewServer()
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
	err = client1.SendMessage("2,3", "msg\n")
	assert.NoError(t, err)

	msgToClient2 := <-client2OutChan
	assert.Equal(t, "from:1 to:2 msg:msg\n", msgToClient2)

	msgToClient3 := <-client3OutChan
	assert.Equal(t, "from:1 to:3 msg:msg\n", msgToClient3)
}
