package binaryxml

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

const malformedErrorStr = "Content is not valid binary XML; %s"

func ToXML(data []byte) (string, error) {
	reader := bytes.NewReader(data)

	// Read table begin marker
	var token BinXMLType
	if err := binary.Read(reader, binary.BigEndian, &token); err != nil {
		return "", err
	}
	if token != tablebegin {
		return "", fmt.Errorf(malformedErrorStr, "missing table begin token")
	}

	// Read table length
	var tableLength uint16
	if err := binary.Read(reader, binary.BigEndian, &tableLength); err != nil {
		return "", err
	}

	// Read table
	elementNamesById := make(map[uint16]string)
	for i := uint16(1); i <= tableLength; i++ {
		name, err := readNullTerminatedString(reader)
		if err != nil {
			return "", err
		}
		elementNamesById[i] = name
	}

	// Read table end marker
	if err := binary.Read(reader, binary.BigEndian, &token); err != nil {
		return "", err
	}
	if token != tableend {
		return "", fmt.Errorf(malformedErrorStr, "missing table end token")
	}

	// Read serial begin marker
	if err := binary.Read(reader, binary.BigEndian, &token); err != nil {
		return "", err
	}
	if token != serialbegin {
		return "", fmt.Errorf(malformedErrorStr, "missing serial begin token")
	}

	// Read serial section
	var xmlBuffer bytes.Buffer
	//xmlBuffer.WriteString("<?xml version=\"1.0\"?>\n")
	if err := readSerialSection(reader, elementNamesById, &xmlBuffer); err != nil {
		return "", err
	}

	return xmlBuffer.String(), nil
}

func readNullTerminatedString(reader io.Reader) (string, error) {
	var buffer bytes.Buffer
	for {
		var byte uint8
		if err := binary.Read(reader, binary.BigEndian, &byte); err != nil {
			return "", err
		}
		if byte == 0 {
			break
		}
		buffer.WriteString(string(byte))
	}
	return buffer.String(), nil
}

func readSerialSection(reader io.Reader, elementNamesById map[uint16]string, response *bytes.Buffer) error {
	elementNameStack := list.New()
	for {
		// Read datatype
		var dataType BinXMLType
		if err := binary.Read(reader, binary.BigEndian, &dataType); err != nil {
			return err
		}

		// Detect serial end marker
		if dataType == serialend {
			return nil
		}

		// Write begin of element
		if isElementType(dataType) {
			var key uint16
			if err := binary.Read(reader, binary.BigEndian, &key); err != nil {
				return err
			}
			elementName, ok := elementNamesById[key]
			if !ok {
				return fmt.Errorf(malformedErrorStr, "no table entry for key")
			}
			elementNameStack.PushFront(elementName)
			response.WriteString("<" + elementName + ">")
		}

		// Write element value
		switch dataType {
		case binarytype:
			var length uint32
			if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
				return err
			}
			value := make([]byte, length)
			if err := binary.Read(reader, binary.BigEndian, &value); err != nil {
				return err
			}
			response.WriteString("[BINARYDATA]")
		case float4type:
			var value float32
			if err := binary.Read(reader, binary.BigEndian, &value); err != nil {
				return err
			}
			response.WriteString(strconv.FormatFloat(float64(value), 'f', 10, 32))
		case int1btype:
			var value int8
			if err := binary.Read(reader, binary.BigEndian, &value); err != nil {
				return err
			}
			response.WriteString(fmt.Sprintf("%d", value))
		case nodetype:
		case uint1btype:
			var value uint8
			if err := binary.Read(reader, binary.BigEndian, &value); err != nil {
				return err
			}
			response.WriteString(fmt.Sprintf("%d", value))
		case int2btype:
			var value int16
			if err := binary.Read(reader, binary.BigEndian, &value); err != nil {
				return err
			}
			response.WriteString(fmt.Sprintf("%d", value))
		case strtype:
			value, err := readNullTerminatedString(reader)
			if err != nil {
				return err
			}
			response.WriteString(value)
		case uint2btype:
			var value uint16
			if err := binary.Read(reader, binary.BigEndian, &value); err != nil {
				return err
			}
			response.WriteString(fmt.Sprintf("%d", value))
		case int4btype:
			var value int32
			if err := binary.Read(reader, binary.BigEndian, &value); err != nil {
				return err
			}
			response.WriteString(fmt.Sprintf("%d", value))
		case uint4btype:
			var value uint32
			if err := binary.Read(reader, binary.BigEndian, &value); err != nil {
				return err
			}
			response.WriteString(fmt.Sprintf("%d", value))
		case int8btype:
			var value int64
			if err := binary.Read(reader, binary.BigEndian, &value); err != nil {
				return err
			}
			response.WriteString(fmt.Sprintf("%d", value))
		case uint8btype:
			var value uint64
			if err := binary.Read(reader, binary.BigEndian, &value); err != nil {
				return err
			}
			response.WriteString(fmt.Sprintf("%d", value))
		}

		// Write end of element
		if dataType == endtagtype {
			element := elementNameStack.Front()
			if element == nil {
				return fmt.Errorf(malformedErrorStr, "too many close element tags")
			}
			elementName := element.Value.(string)
			elementNameStack.Remove(element)
			response.WriteString("</" + elementName + ">")
		}
	}
	return nil
}

func isElementType(x BinXMLType) bool {
	if x == nodetype || x == int1btype || x == int2btype || x == int4btype || x == int8btype {
		return true
	}
	if x == uint1btype || x == uint2btype || x == uint4btype || x == uint8btype {
		return true
	}
	if x == float4type || x == strtype || x == binarytype {
		return true
	}
	return false
}
