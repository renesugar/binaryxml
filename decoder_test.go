package binaryxml_test

import (
	"fmt"
	"github.com/BixData/binaryxml"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/xml"
	"io/ioutil"
	"testing"
)


func TestDecodeBinaryFixture1(t *testing.T) {
	doDecodeTest("test-systemlib-1", t)
}

func TestDecodeBinaryFixture2(t *testing.T) {
	doDecodeTest("test-systemlib-2", t)
}

func TestDecodeBinaryFixture3(t *testing.T) {
	doDecodeTest("test-systemlib-3", t)
}

func TestDecodeBinaryFixture4(t *testing.T) {
	doDecodeTest("test-systemlib-4", t)
}

func TestDecodeBinaryFixture5(t *testing.T) {
	doDecodeTest("test-systemlib-5", t)
}


func doDecodeTest(fixtureName string, t *testing.T) {
	// Configure XML minifier, and use during comparisons to minimize superficial differences
	minifier := minify.New()
	minifier.AddFunc("text/xml", xml.Minify) 

	// Read fixtures
	fixture := "testdata/" + fixtureName + ".binaryxml"
	fmt.Printf("Testing decode with %s\n", fixture)
	binaryXml, err := ioutil.ReadFile(fixture)
	if err != nil {
		t.Errorf("Failed opening %s", fixture)
	}
	fixture = "testdata/" + fixtureName + ".xml"
	expectedXmlBinary, err := ioutil.ReadFile(fixture)
	if err != nil {
		t.Errorf("Failed opening %s", fixture)
	}
	expectedXml := string(expectedXmlBinary)
	
	// Decode binary xml
	xml, err := binaryxml.Decode(binaryXml)
	if err != nil {
		t.Errorf("Failed decoding binary xml %+v", err)
	}
	
	// Minify expected vs actual xml, to minimize inconsequential differences
	xml, err = minifier.String("text/xml", xml)
	if err != nil {
		t.Errorf("Failed minifying xml for comparison %+v", err)
	}
	var minifiedXml string
	var minifiedExpectedXml string
	minifiedXml, err = minifier.String("text/xml", xml)
	minifiedExpectedXml, err = minifier.String("text/xml", expectedXml)
	
	// compare
	if (minifiedXml != minifiedExpectedXml) {
		t.Errorf("Failed converting binary xml; expected %s; got %s", expectedXml, xml)
	}
}
