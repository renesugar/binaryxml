package binaryxml

import (
	"io"
)


func Encode(value interface{}, writer io.Writer) error {
	encoder := NewEncoder(writer)
	return encoder.Encode(value);
}
