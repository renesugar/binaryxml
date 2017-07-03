package binaryxml_test

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/BixData/binaryxml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)


type Fixture1 struct {
	XMLName     struct{} `xml:"BixRequest"`
	Request     string   `xml:"request"`
	ToNamespace string   `xml:"toNamespace"`
	MOID        string   `xml:"moid"`
	MID         string   `xml:"mid"`
}


type Fixture4 struct {
	XMLName   struct{} `xml:"TestDoc"`
	Int8Min   int8     `xml:"int8_min"`
	Int8Max   int8     `xml:"int8_max"`
	Int16Min  int16    `xml:"int16_min"`
	Int16Max  int16    `xml:"int16_max"`
	Int32Min  int32    `xml:"int32_min"`
	Int32Max  int32    `xml:"int32_max"`
	Int64Min  int64    `xml:"int64_min"`
	Int64Max  int64    `xml:"int64_max"`
	Uint8Max  uint8    `xml:"uint8_max"`
	Uint16Max uint16   `xml:"uint16_max"`
	Uint32Max uint32   `xml:"uint32_max"`
	Uint64Max uint64   `xml:"uint64_max"`
}
	

func TestEncodeFixture1(t *testing.T) {
	assert := assert.New(t)
	
	// Read XML fixture #1
	fixture := "testdata/test-systemlib-1.xml"
	fmt.Printf("Loading fixture %s\n", fixture)
	xmlBytes, err := ioutil.ReadFile(fixture)
	if err != nil {t.Errorf("Failed opening test fixture %s", fixture)}
	
	// Unmarshal XML file into Fixture1 structure
	fixture1 := Fixture1{}
	err = xml.Unmarshal(xmlBytes, &fixture1)
	if err != nil {t.Errorf("Failed unmarshalling test fixture #1 into structure: %v", err)}
	
	// Sanity check - ensure bixRequest contains expected values loaded from file, before starting test
	assert.Equal("VirtualMachines", fixture1.ToNamespace, "Failed loading %s into struct %+v", fixture, fixture1)
	assert.Equal("Testing", fixture1.Request, "Failed loading %s into struct %+v", fixture, fixture1)
	assert.Equal("6", fixture1.MOID, "Failed loading %s into struct %+v", fixture, fixture1)
	assert.Equal("1", fixture1.MID, "Failed loading %s into struct %+v", fixture, fixture1)
	
	// Encode structure as binary xml
	file, _ := ioutil.TempFile("", "binaryxmlEncoderTest1")
	writer := bufio.NewWriter(file)
	if err := binaryxml.Encode(fixture1, writer); err != nil {t.Errorf("Failed encoding object as binary xml %+v", err)}
	writer.Flush()
	
	// Unmarshal binary XML file into 2nd Fixture1 structure
	binaryXmlBytes, err := ioutil.ReadFile(file.Name())
	if err != nil {t.Errorf("Failed opening generated binary xml file %s", file.Name)}
	xmlString, err := binaryxml.Decode(binaryXmlBytes)
	if err != nil {t.Errorf("Failed decoding binary xml %+v", err)}	
	secondFixture1 := Fixture1{}
	err = xml.Unmarshal([]byte(xmlString), &secondFixture1)
	if err != nil {t.Errorf("Failed unmarshalling test fixture #1 into structure: %v", err)}
	assert.Equal(fixture1, secondFixture1)
}


func TestEncodeFixture4(t *testing.T) {
	assert := assert.New(t)
	
	// Read XML fixture #4
	fixture := "testdata/test-systemlib-4.xml"
	fmt.Printf("Loading fixture %s\n", fixture)
	xmlBytes, err := ioutil.ReadFile(fixture)
	if err != nil {t.Errorf("Failed opening test fixture %s", fixture)}
	
	// Unmarshal XML file into Fixture4 structure
	fixture4 := Fixture4{}
	err = xml.Unmarshal(xmlBytes, &fixture4)
	if err != nil {t.Errorf("Failed unmarshalling test fixture #4 into structure: %v", err)}
	
	// Sanity check - ensure bixRequest contains expected values loaded from file, before starting test
	assert.Equal(int8(-128), fixture4.Int8Min  , "Failed loading %s into struct %+v", fixture, fixture4)
	assert.Equal(int8(127), fixture4.Int8Max  , "Failed loading %s into struct %+v", fixture, fixture4)
	assert.Equal(int16(-32768), fixture4.Int16Min , "Failed loading %s into struct %+v", fixture, fixture4)
	assert.Equal(int16(32767), fixture4.Int16Max , "Failed loading %s into struct %+v", fixture, fixture4)
	assert.Equal(int32(-2147483648), fixture4.Int32Min , "Failed loading %s into struct %+v", fixture, fixture4)
	assert.Equal(int32(2147483647), fixture4.Int32Max , "Failed loading %s into struct %+v", fixture, fixture4)
	assert.Equal(int64(-9223372036854775808), fixture4.Int64Min , "Failed loading %s into struct %+v", fixture, fixture4)
	assert.Equal(int64(9223372036854775807), fixture4.Int64Max , "Failed loading %s into struct %+v", fixture, fixture4)
	assert.Equal(uint8(255), fixture4.Uint8Max , "Failed loading %s into struct %+v", fixture, fixture4)
	assert.Equal(uint16(65535), fixture4.Uint16Max, "Failed loading %s into struct %+v", fixture, fixture4)
	assert.Equal(uint32(4294967295), fixture4.Uint32Max, "Failed loading %s into struct %+v", fixture, fixture4)
	assert.Equal(uint64(18446744073709551615), fixture4.Uint64Max, "Failed loading %s into struct %+v", fixture, fixture4)
	
	// Encode structure as binary xml
	file, _ := ioutil.TempFile("", "binaryxmlEncoderTest4")
	writer := bufio.NewWriter(file)
	fmt.Printf("Writing binary encoded xml file %s\n", file.Name())
	if err := binaryxml.Encode(fixture4, writer); err != nil {t.Errorf("Failed encoding object as binary xml %+v", err)}
	writer.Flush()
	
	// Unmarshal binary XML file into 2nd Fixture4 structure
	binaryXmlBytes, err := ioutil.ReadFile(file.Name())
	if err != nil {t.Errorf("Failed opening generated binary xml file %s", file.Name)}
	xmlString, err := binaryxml.Decode(binaryXmlBytes)
	if err != nil {t.Errorf("Failed decoding binary xml %+v", err)}	
	secondFixture4 := Fixture4{}
	err = xml.Unmarshal([]byte(xmlString), &secondFixture4)
	if err != nil {t.Errorf("Failed unmarshalling test fixture #4 into structure: %v", err)}
	assert.Equal(fixture4, secondFixture4)
}
