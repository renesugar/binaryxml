package client

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	assert := assert.New(t)

	listener, err := net.Listen("tcp", "localhost:0")
	assert.NoError(err)
	port := listener.Addr().(*net.TCPAddr).Port

	accepted := false
	go func() {
		if _, err := listener.Accept(); err == nil {
			accepted = true
		}
	}()

	_, err = Connect("127.0.0.1", port)
	assert.NoError(err)

	time.Sleep(20 * time.Millisecond)
	assert.True(accepted)
	assert.NoError(listener.Close())
}
