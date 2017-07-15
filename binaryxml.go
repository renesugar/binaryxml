package binaryxml

import ()

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

type BixError struct {
	XMLName       struct{} `xml:"BixError"`
	FromNamespace string   `xml:"fromNamespace"`
	Request       string   `xml:"request"`
	MOID          uint64   `xml:"moid"`
	MID           uint64   `xml:"mid"`
	Error         string   `xml:"error"`
}
