package server

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_ShouldStart(t *testing.T) {
	s := NewServer(&sync.Mutex{})
	err := s.Start("1234")
	defer s.Close()

	assert.NoError(t, err)
}
