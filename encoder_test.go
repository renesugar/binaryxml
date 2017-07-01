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


type BixRequest struct {
	XMLName     struct{} `xml:"BixRequest"`
	Request     string   `xml:"request"`
	ToNamespace string   `xml:"toNamespace"`
	MOID        uint64   `xml:"moid"`
	MID         uint64   `xml:"mid"`
}


func TestEncodeFixture1(t *testing.T) {
	assert := assert.New(t)
	
	// Read XML fixture #1
	fixture := "testdata/test-systemlib-1.xml"
	fmt.Printf("Loading fixture %s\n", fixture)
	xmlBytes, err := ioutil.ReadFile(fixture)
	if err != nil {t.Errorf("Failed opening test fixture %s", fixture)}
	
	// Unmarshal XML file into BixRequest structure
	bixRequest := BixRequest{}
	err = xml.Unmarshal(xmlBytes, &bixRequest)
	if err != nil {t.Errorf("Failed unmarshalling test fixture #1 into structure: %v", err)}
	
	// Sanity check - ensure bixRequest contains expected values loaded from file, before starting test
	assert.Equal(bixRequest.ToNamespace, "VirtualMachines", "Failed loading %s into struct %+v", fixture, bixRequest)
	assert.Equal(bixRequest.Request, "Testing", "Failed loading %s into struct %+v", fixture, bixRequest)
	assert.Equal(bixRequest.MOID, uint64(6), "Failed loading %s into struct %+v", fixture, bixRequest)
	assert.Equal(bixRequest.MID, uint64(1), "Failed loading %s into struct %+v", fixture, bixRequest)
	
	// Encode structure as binary xml
	file, _ := ioutil.TempFile("", "binaryxmlEncoderTest")
	writer := bufio.NewWriter(file)
	if err := binaryxml.Encode(bixRequest, writer); err != nil {t.Errorf("Failed encoding object as binary xml %+v", err)}
	writer.Flush()
	
	// Compare
	fixture = "testdata/test-systemlib-1.binaryxml"
	binaryXml, err := ioutil.ReadFile(file.Name())
	if err != nil {t.Errorf("Failed opening generated binary xml file %s", file.Name)}
	expectedBinaryXml, err := ioutil.ReadFile(fixture)
	if err != nil {t.Errorf("Failed opening test fixture %s", fixture)}
	assert.Equal(expectedBinaryXml, binaryXml, "Binary xml does not match")
}
