package binaryxml_client

import (
	"bufio"
	"bytes"
	"fmt"
	"net"

	"github.com/BixData/binaryxml"
	"github.com/BixData/binaryxml/messages"
	"github.com/docktermj/go-logger/logger"
)

// ----------------------------------------------------------------------------

type Client struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func (self *Client) Close() error {
	return self.Close()
}

func (self *Client) SendRaw(param uint8, binaryXML []byte) error {
	if err := binaryxml_messages.WriteMessage(self.Writer, param, binaryXML); err != nil {
		return err
	}
	return self.Writer.Flush()
}

func (self *Client) Send(param uint8, req interface{}) error {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	if err := binaryxml.Encode(req, writer); err != nil {
		return err
	}
	writer.Flush()
	binaryXML := buffer.Bytes()
	return self.SendRaw(param, binaryXML)
}

func (self *Client) ReceiveRaw(param *uint8, binaryXML *[]byte) error {
	return binaryxml_messages.ReadMessage(self.Reader, param, binaryXML)
}

// ----------------------------------------------------------------------------

func Connect(host string, port int) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		logger.Warnf("Failed connecting to %s: %v", addr, err)
		return nil, err
	}
	logger.Debugf("Connected to %s", addr)
	client := Client{Conn: conn, Reader: bufio.NewReader(conn), Writer: bufio.NewWriter(conn)}
	return &client, nil
}
