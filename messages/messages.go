package binaryxml_messages

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
)

const (
	msgstate_start uint8 = 121
	msgstate_end   uint8 = 123
)

// ----------------------------------------------------------------------------
// Reads a message
// ----------------------------------------------------------------------------

func ReadMessage(reader io.Reader, param *uint8, binaryXML *[]byte) error {
	// Read message start token
	var token uint8
	if err := binary.Read(reader, binary.BigEndian, &token); err != nil {
		return err
	}
	if token != msgstate_start {
		return fmt.Errorf("Malformed message; missing start token")
	}

	// Read message length
	var messageLength uint32
	if err := binary.Read(reader, binary.BigEndian, &messageLength); err != nil {
		return err
	}

	// Read param
	if err := binary.Read(reader, binary.BigEndian, param); err != nil {
		return err
	}

	// Read message
	*binaryXML = make([]byte, messageLength)
	if _, err := io.ReadFull(reader, *binaryXML); err != nil {
		return err
	}

	// Read message end token
	if err := binary.Read(reader, binary.BigEndian, &token); err != nil {
		return err
	}
	if token != msgstate_end {
		return fmt.Errorf("Malformed message; missing end token")
	}

	// Read crc32 checksum
	var crcFromPayload uint32
	if err := binary.Read(reader, binary.BigEndian, &crcFromPayload); err != nil {
		return err
	}

	// Compute our own crc32 for the message, and reject message whose checksum doesn't match
	locallyComputedCrc := crc32.ChecksumIEEE(*binaryXML)
	if locallyComputedCrc != crcFromPayload {
		return fmt.Errorf("Malformed message; crc32 checksum does not match")
	}

	return nil
}

func WriteMessage(writer io.Writer, param uint8, binaryXML []byte) error {
	// Write message start token
	if err := binary.Write(writer, binary.BigEndian, msgstate_start); err != nil {
		return err
	}

	// Write message length
	if err := binary.Write(writer, binary.BigEndian, uint32(len(binaryXML))); err != nil {
		return err
	}

	// Write param
	if err := binary.Write(writer, binary.BigEndian, param); err != nil {
		return err
	}

	// Write message
	if _, err := writer.Write(binaryXML); err != nil {
		return err
	}

	// Write message end token
	if err := binary.Write(writer, binary.BigEndian, msgstate_end); err != nil {
		return err
	}

	// Write crc32 checksum
	crc := crc32.ChecksumIEEE(binaryXML)
	if err := binary.Write(writer, binary.BigEndian, crc); err != nil {
		return err
	}

	return nil
}
