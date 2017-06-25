package binaryxml

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const (
)

type BinXMLType uint8
const (
	undefinedtype BinXMLType = iota
	nodetype
	int1btype
	uint1btype
	int2btype
	uint2btype
	int4btype
	uint4btype
	int8btype
	uint8btype
	float4type
	strtype
	binarytype
	endtagtype
	
	tablebegin  BinXMLType = 124
	tableend    BinXMLType = 125
	serialbegin BinXMLType = 126
	serialend   BinXMLType = 127
)

var malformedError = errors.New("Content is not valid binary XML")


func Decode(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	
	// Read table begin marker
	var token BinXMLType
	err := binary.Read(reader, binary.BigEndian, &token)
	if err != nil || token != tablebegin {return "", malformedError}
	
	// Read table length
	var tableLength uint16
	err = binary.Read(reader, binary.BigEndian, &tableLength)
	if err != nil {return "", err}
	
	// Read table
	elementNamesById := make(map[uint16]string)
	for i := uint16(1); i <= tableLength; i++ {
		name, err := readNullTerminatedString(reader)
		if err != nil {return "", err}
		elementNamesById[i] = name
	}
// 	fmt.Printf("elementNamesById: %v\n", elementNamesById)
	
	// Read table end marker
	err = binary.Read(reader, binary.BigEndian, &token)
	if err != nil || token != tableend {return "", malformedError}
	
	// Read serial begin marker
	err = binary.Read(reader, binary.BigEndian, &token)
	if err != nil || token != serialbegin {return "", malformedError}
	
	// Read serial section
	var xmlBuffer bytes.Buffer
	xmlBuffer.WriteString("<?xml version=\"1.0\"?>\n")
	err = readSerialSection(reader, elementNamesById, &xmlBuffer)
	if err != nil {return "", malformedError}
	
	return xmlBuffer.String(), nil
}


func readNullTerminatedString(reader io.Reader) (string, error) {
	var buffer bytes.Buffer
	for {
		var byte uint8
		if err := binary.Read(reader, binary.BigEndian, &byte); err != nil {return "", malformedError}
		if byte == 0 {break}
		buffer.WriteString(string(byte))
	}
	return buffer.String(), nil
}


func readSerialSection(reader io.Reader, elementNamesById map[uint16]string, response* bytes.Buffer) error {
	elementNameStack := list.New()
	for {
		// Read datatype
		var dataType BinXMLType
		err := binary.Read(reader, binary.BigEndian, &dataType)
		if err != nil {return err}
		
		// Detect serial end marker
		if dataType == serialend {return nil}
		
		// Write begin of element
		if isElementType(dataType) {
			var key uint16
			if err = binary.Read(reader, binary.BigEndian, &key); err != nil {return malformedError}
			elementName, ok := elementNamesById[key]
			if !ok {return malformedError}
			elementNameStack.PushFront(elementName)
			response.WriteString("<" + elementName + ">")
		}
		
		// Write element value
		switch dataType {
		case float4type:
			var value float32
			if err = binary.Read(reader, binary.BigEndian, &value); err != nil {return malformedError}
			response.WriteString(strconv.FormatFloat(float64(value), 'f', 6, 32))
		case int1btype:
			var value int8
			if err = binary.Read(reader, binary.BigEndian, &value); err != nil {return malformedError}
			response.WriteString(fmt.Sprintf("%d", value))
		case nodetype:
		case uint1btype:
			var value uint8
			if err = binary.Read(reader, binary.BigEndian, &value); err != nil {return malformedError}
			response.WriteString(fmt.Sprintf("%d", value))
		case int2btype:
			var value int16
			if err = binary.Read(reader, binary.BigEndian, &value); err != nil {return malformedError}
			response.WriteString(fmt.Sprintf("%d", value))
		case strtype:
			value, err := readNullTerminatedString(reader)
			if err != nil {return malformedError}
			response.WriteString(value)
		case uint2btype:
			var value uint16
			if err = binary.Read(reader, binary.BigEndian, &value); err != nil {return malformedError}
			response.WriteString(fmt.Sprintf("%d", value))
		case int4btype:
			var value int32
			if err = binary.Read(reader, binary.BigEndian, &value); err != nil {return malformedError}
			response.WriteString(fmt.Sprintf("%d", value))
		case uint4btype:
			var value uint32
			if err = binary.Read(reader, binary.BigEndian, &value); err != nil {return malformedError}
			response.WriteString(fmt.Sprintf("%d", value))
		case int8btype:
			var value int64
			if err = binary.Read(reader, binary.BigEndian, &value); err != nil {return malformedError}
			response.WriteString(fmt.Sprintf("%d", value))
		case uint8btype:
			var value uint64
			if err = binary.Read(reader, binary.BigEndian, &value); err != nil {return malformedError}
			response.WriteString(fmt.Sprintf("%d", value))
		}
		
		// Write end of element
		if dataType == endtagtype {
			element := elementNameStack.Front()
			if element == nil {return malformedError}
			elementName := element.Value.(string)
			elementNameStack.Remove(element)
			response.WriteString("</" + elementName + ">")
		}
	}
	return nil
}


func isElementType(x BinXMLType) bool {
	if x==nodetype || x==int1btype || x==int2btype || x==int4btype || x==int8btype {return true}
	if x==uint1btype || x==uint2btype || x==uint4btype || x==uint8btype {return true}
	if x==float4type || x==strtype || x==binarytype {return true}
	return false
}
