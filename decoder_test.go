package binaryxml_test

import (
	"fmt"
	"github.com/BixData/binaryxml"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/xml"
	"io/ioutil"
	"testing"
)


func TestReadBinaryFixture1(t *testing.T) {
	doTest("test-systemlib-1", t)
}

func TestReadBinaryFixture2(t *testing.T) {
	doTest("test-systemlib-2", t)
}

func TestReadBinaryFixture3(t *testing.T) {
	doTest("test-systemlib-3", t)
}

func TestReadBinaryFixture4(t *testing.T) {
	doTest("test-systemlib-4", t)
}


func doTest(fixtureName string, t *testing.T) {
	// Configure XML minifier, and use during comparisons to minimize superficial differences
	minifier := minify.New()
	minifier.AddFunc("text/xml", xml.Minify) 

	// Read fixtures
	fixture := "testdata/" + fixtureName + ".binaryxml"
	fmt.Printf("Testing fixture %s\n", fixture)
	binaryXml, err := ioutil.ReadFile(fixture)
	if err != nil {
		t.Errorf("Failed opening test fixture %s", fixture)
	}
	fixture = "testdata/" + fixtureName + ".xml"
	expectedXmlBinary, err := ioutil.ReadFile(fixture)
	if err != nil {
		t.Errorf("Failed opening test fixture %s", fixture)
	}
	expectedXml := string(expectedXmlBinary)
	
	// Decode binary xml
	xml, err := binaryxml.Decode(binaryXml)
	if err != nil {
		t.Errorf("Failed decoding binary xml %+v", err)
	}
	xml, err = minifier.String("text/xml", xml)
	if err != nil {
		t.Errorf("Failed minifying xml for comparison %+v", err)
	}
	
	// Minify expected vs actual xml, then compare
	var minifiedXml string
	var minifiedExpectedXml string
	minifiedXml, err = minifier.String("text/xml", xml)
	minifiedExpectedXml, err = minifier.String("text/xml", expectedXml)
	if (minifiedXml != minifiedExpectedXml) {
		t.Errorf("Failed converting binary xml; expected %s; got %s", expectedXml, xml)
	}
}
