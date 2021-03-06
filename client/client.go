package client

import (
	"bufio"
	"bytes"
	"errors"
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
	if err := messages.WriteMessage(self.Writer, param, binaryXML); err != nil {
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
	return messages.ReadMessage(self.Reader, param, binaryXML)
}

func (self *Client) Receive(param *uint8, res interface{}) error {
	var binaryXML []byte
	if err := self.ReceiveRaw(param, &binaryXML); err != nil {
		return err
	}
	err := binaryxml.Decode(binaryXML, &res)
	if err != nil {
		var bixError binaryxml.BixError
		err2 := binaryxml.Decode(binaryXML, &bixError)
		if err2 == nil && bixError.Error != "" {
			return errors.New(bixError.Error)
		}
	}
	return err
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
