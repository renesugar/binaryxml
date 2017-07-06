package binaryxml

import (
	"io"
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


func NewEncoder(writer io.Writer) *BinaryXMLEncoder {
	return &BinaryXMLEncoder{writer:writer}
}
