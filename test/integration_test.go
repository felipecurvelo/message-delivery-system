package test

import (
	"testing"

	"github.com/felipecurvelo/message-delivery-system/internal/client"
	"github.com/felipecurvelo/message-delivery-system/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	s := server.NewServer()
	s.Start("1234")
	defer s.Close()

	client1 := client.NewClient()
	err := client1.Connect("127.0.0.1:1234")
	assert.NoError(t, err)

	client2 := client.NewClient()
	err = client2.Connect("127.0.0.1:1234")
	assert.NoError(t, err)

	client1OutChan := make(chan string)
	go client1.HandleMessages(client1OutChan)

	client2OutChan := make(chan string)
	go client2.HandleMessages(client2OutChan)

	//Send a first msg
	err = client1.SendMessage("2", "firstMsgToClient2\n")
	assert.NoError(t, err)

	err = client2.SendMessage("1", "firstMsgToClient1\n")
	assert.NoError(t, err)

	firstMsgToClient1 := <-client1OutChan
	assert.Equal(t, "from:2 to:1 msg:firstMsgToClient1\n", firstMsgToClient1)

	firstMsgToClient2 := <-client2OutChan
	assert.Equal(t, "from:1 to:2 msg:firstMsgToClient2\n", firstMsgToClient2)

	//Send a second msg
	err = client1.SendMessage("2", "secondMsgToClient2\n")
	assert.NoError(t, err)

	err = client2.SendMessage("1", "secondMsgToClient1\n")
	assert.NoError(t, err)

	secondMsgToClient1 := <-client1OutChan
	assert.Equal(t, "from:2 to:1 msg:secondMsgToClient1\n", secondMsgToClient1)

	secondMsgToClient2 := <-client2OutChan
	assert.Equal(t, "from:1 to:2 msg:secondMsgToClient2\n", secondMsgToClient2)
}
