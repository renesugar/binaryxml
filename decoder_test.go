package binaryxml_test

import (
	"fmt"
	"github.com/BixData/binaryxml"
	"io/ioutil"
	"testing"
)


func TestReadBinaryFixture1(t *testing.T) {
	// Read fixtures
	fixture := "testdata/test-systemlib-1.binaryxml"
	binaryXml, err := ioutil.ReadFile(fixture)
	if err != nil {
		t.Errorf("Failed opening test fixture %s", fixture)
	}
	fixture = "testdata/test-systemlib-1.xml"
	expectedXmlBinary, err := ioutil.ReadFile(fixture)
	if err != nil {
		t.Errorf("Failed opening test fixture %s", fixture)
	}
	expectedXml := string(expectedXmlBinary)
	
	// Decode binary xml and test
	xml := binaryxml.Decode(binaryXml)
	fmt.Println(xml)
	if (xml != expectedXml) {
		t.Errorf("Failed converting binary xml; expected %s; got %s", expectedXml, xml)
	}
}
