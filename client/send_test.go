package client

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/BixData/binaryxml/messages"
	"github.com/docktermj/go-logger/logger"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	logger.SetLevel(logger.LevelDebug)
	assert := assert.New(t)

	// Create a server
	listener, err := net.Listen("tcp", "localhost:0")
	assert.NoError(err)
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	logger.Debugf("Listening on port %d", port)

	var conn net.Conn
	accepted := false
	received := false
	var listenerErr error
	go func() {
		if conn, listenerErr = listener.Accept(); listenerErr != nil {
			return
		}
		defer conn.Close()
		accepted = true

		// Read message
		reader := bufio.NewReader(conn)
		var param uint8
		var binaryXML []byte
		if listenerErr = messages.ReadMessage(reader, &param, &binaryXML); listenerErr != nil {
			return
		}
		received = true

	}()

	// Create a client
	client, err := Connect("127.0.0.1", port)
	assert.NoError(err)

	time.Sleep(20 * time.Millisecond)
	assert.True(accepted)

	// Prepare a BinaryXML request
	type MyRequest struct {
		XMLName     struct{} `xml:"BixRequest"`
		Request     string   `xml:"request"`
		ToNamespace string   `xml:"toNamespace"`
	}
	myRequest := MyRequest{ToNamespace: "foo", Request: "bar"}

	// Send it
	assert.NoError(client.Send(0, myRequest))

	time.Sleep(200 * time.Millisecond)
	assert.True(received)
}
