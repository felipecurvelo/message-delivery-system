package client

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_WhenMessageIsEmpty_ShouldReturnError(t *testing.T) {
	client := NewClient()
	err := client.SendMessage("1", []byte(""))
	assert.Error(t, err)
	assert.Equal(t, "Message is empty", err.Error())
}

func TestClient_WhenDestinationIsEmpty_ShouldReturnError(t *testing.T) {
	client := NewClient()
	err := client.SendMessage("", []byte("msg"))
	assert.Error(t, err)
	assert.Equal(t, "Destination is empty", err.Error())
}

func TestClient_WhenMoreThanMaxMessageSize_ShouldReturnError(t *testing.T) {
	msg := make([]byte, 1024001)
	client := NewClient()
	err := client.SendMessage("1", msg)
	assert.Error(t, err)
	assert.Equal(t, "Message size should be less than 1024Kb", err.Error())
}

func TestClient_WhenMoreThanMaxDestination_ShouldReturnError(t *testing.T) {
	destinationIDs := []string{}
	for i := 0; i <= 300; i++ {
		destinationIDs = append(destinationIDs, strconv.Itoa(i))
	}

	client := NewClient()
	err := client.SendMessage(strings.Join(destinationIDs, ","), []byte("msg"))
	assert.Error(t, err)
	assert.Equal(t, "Max destination per message reached", err.Error())
}
