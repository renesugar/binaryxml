package binaryxml

import (
	"encoding/xml"
)

func Decode(binaryXML []byte, v interface{}) error {
	xmlString, err := ToXML(binaryXML)
	if err != nil {
		return err
	}
	return xml.Unmarshal([]byte(xmlString), v)
}
