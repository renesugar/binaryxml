package binaryxml_client_test

import (
	"bufio"
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/BixData/binaryxml"
	"github.com/BixData/binaryxml/client"
	"github.com/BixData/binaryxml/messages"
	"github.com/docktermj/go-logger/logger"
	"github.com/stretchr/testify/assert"
)

func TestReceiveError(t *testing.T) {
	logger.SetLevel(logger.LevelDebug)
	assert := assert.New(t)

	type MyRequest struct {
		XMLName     struct{} `xml:"BixRequest"`
		Request     string   `xml:"request"`
		ToNamespace string   `xml:"toNamespace"`
	}

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
			logger.Errorf("%v", listenerErr)
			return
		}
		defer conn.Close()
		accepted = true

		// Read message
		reader := bufio.NewReader(conn)
		var param uint8
		var binaryXML []byte
		if listenerErr = binaryxml_messages.ReadMessage(reader, &param, &binaryXML); listenerErr != nil {
			logger.Errorf("%v", listenerErr)
			return
		}
		received = true

		// Prepare response
		myRes := binaryxml.BixError{FromNamespace: "baz", Error: "567"}
		var buffer bytes.Buffer
		{
			writer := bufio.NewWriter(&buffer)
			if listenerErr = binaryxml.Encode(myRes, writer); listenerErr != nil {
				logger.Errorf("%v", listenerErr)
				return
			}
			writer.Flush()
		}
		binaryXML = buffer.Bytes()

		// Send response
		writer := bufio.NewWriter(conn)
		if listenerErr = binaryxml_messages.WriteMessage(writer, param, binaryXML); listenerErr != nil {
			logger.Errorf("%v", listenerErr)
			return
		}
		writer.Flush()
	}()

	// Create a client
	client, err := binaryxml_client.Connect("127.0.0.1", port)
	assert.NoError(err)

	time.Sleep(20 * time.Millisecond)
	assert.True(accepted)

	// Prepare and send request
	myRequest := MyRequest{ToNamespace: "foo", Request: "bar"}
	assert.NoError(client.Send(0, myRequest))
	time.Sleep(200 * time.Millisecond)
	assert.True(received)

	// Receive response
	var param uint8

	type MyResponse struct {
		XMLName       struct{} `xml:"BixResponse"`
		FromNamespace string   `xml:"fromNamespace"`
	}
	var myRes MyResponse
	err = client.Receive(&param, &myRes)
	assert.Error(err)
	assert.Equal("567", err.Error())
}
