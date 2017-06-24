package binaryxml

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"fmt"
)

const tablebegin  uint8 = 124
const tableend    uint8 = 125
const serialbegin uint8 = 126
const serialend   uint8 = 127


func Decode(data []byte) (string, error) {
	fmt.Println("Decode()")
	
	reader := bytes.NewReader(data)
	var uint8Value uint8
	err := binary.Read(reader, binary.BigEndian, &uint8Value)
	if err != nil || uint8Value != tablebegin {
		return "", errors.New("Content is not binary XML")
	}
	
	// Read table length
	var tableLength uint16
	err = binary.Read(reader, binary.BigEndian, &tableLength)
	if err != nil {return "", err}
	fmt.Printf("Table length: %d\n", tableLength)
	
	// Read table
	elementNamesById := make(map[uint16]string)
	for i := uint16(1); i <= tableLength; i++ {
		name, err := readNullTerminatedString(reader)
		if err != nil {return "", err}
		elementNamesById[i] = name
	}
	fmt.Printf("elementNamesById: %v\n", elementNamesById)
	
	return "Foo", nil // TODO
}


func readNullTerminatedString(reader io.ByteReader) (string, error) {
	var buffer bytes.Buffer
	for {
		byte, err := reader.ReadByte()
		if err != nil || byte == 0 {break}
		buffer.WriteString(string(byte))
	}
	return buffer.String(), nil
}
