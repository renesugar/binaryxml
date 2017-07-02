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
	MOID        string   `xml:"moid"`
	MID         string   `xml:"mid"`
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
	assert.Equal(bixRequest.MOID, "6", "Failed loading %s into struct %+v", fixture, bixRequest)
	assert.Equal(bixRequest.MID, "1", "Failed loading %s into struct %+v", fixture, bixRequest)
	
	// Encode structure as binary xml
	file, _ := ioutil.TempFile("", "binaryxmlEncoderTest")
	writer := bufio.NewWriter(file)
	if err := binaryxml.Encode(bixRequest, writer); err != nil {t.Errorf("Failed encoding object as binary xml %+v", err)}
	writer.Flush()
	
	// Unmarshal binary XML file into 2nd BixRequest structure
	binaryXmlBytes, err := ioutil.ReadFile(file.Name())
	if err != nil {t.Errorf("Failed opening generated binary xml file %s", file.Name)}
	xmlString, err := binaryxml.Decode(binaryXmlBytes)
	if err != nil {t.Errorf("Failed decoding binary xml %+v", err)}	
	bixRequest2 := BixRequest{}
	err = xml.Unmarshal([]byte(xmlString), &bixRequest2)
	if err != nil {t.Errorf("Failed unmarshalling test fixture #1 into structure: %v", err)}
}
