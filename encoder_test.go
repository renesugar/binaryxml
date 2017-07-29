package binaryxml_test

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/BixData/binaryxml"
	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify"
	xmlminifier "github.com/tdewolff/minify/xml"
	"io/ioutil"
	"testing"
)

// Numbered fixtures below correspond to C++ testdata. Alpha fixtures are original to this file.

type Fixture1 struct {
	XMLName     struct{} `xml:"BixRequest"`
	Request     string   `xml:"request"`
	ToNamespace string   `xml:"toNamespace"`
	MOID        string   `xml:"moid"`
	MID         string   `xml:"mid"`
}

type Fixture2 struct {
	XMLName     struct{}      `xml:"BixRequest"`
	ToNamespace string        `xml:"toNamespace"`
	Request     string        `xml:"request"`
	MOID        string        `xml:"moid"`
	MID         string        `xml:"mid"`
	Data        Fixture2_Data `xml:"Data"`
}

type Fixture2_Data struct {
	XMLName struct{}       `xml:"Data"`
	MOID    string         `xml:"moid"`
	Query   Fixture2_Query `xml:"Query"`
}

type Fixture2_Query struct {
	XMLName    struct{} `xml:"Query"`
	Namespace  string   `xml:"namespace"`
	Instance   string   `xml:"instance"`
	Key        string   `xml:"key"`
	Field      string   `xml:"field"`
	Interval   uint32   `xml:"interval"`
	UserKey    string   `xml:"userkey"`
	Tablespace string   `xml:"tablespace"`
}

type Fixture3 struct {
	XMLName   struct{}           `xml:"DataPluginDefinition"`
	Namespace Fixture3_Namespace `xml:"Namespace"`
}

type Fixture3_Namespace struct {
	XMLName struct{}        `xml:"Namespace"`
	Name    string          `xml:"name"`
	Display string          `xml:"display"`
	Version string          `xml:"version"`
	Schema  Fixture3_Schema `xml:"Schema"`
}

type Fixture3_Schema struct {
	XMLName  struct{}          `xml:"Schema"`
	Name     string            `xml:"name"`
	Instance Fixture3_Instance `xml:"Instance"`
	Fields   []Fixture3_Field  `xml:"Field"`
}

type Fixture3_Instance struct {
	XMLName struct{} `xml:"Instance"`
	Name    string   `xml:"name"`
	Display string   `xml:"display"`
}

type Fixture3_Field struct {
	XMLName  struct{} `xml:"Field"`
	Name     string   `xml:"name"`
	Type     string   `xml:"type"`
	Display  string   `xml:"display"`
	UserData string   `xml:"userdata"`
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

type Fixture5 struct {
	XMLName            struct{} `xml:"TestDoc"`
	Float32_0          float32  `xml:"float32_0"`
	Float32_Pi         float32  `xml:"float32_pi"`
	Float32_NegativePi float32  `xml:"float32_negativepi"`
}

type Fixture6 struct {
	XMLName         struct{} `xml:"TestDoc"`
	Binary_007f80ff []byte   `xml:"binary_007f80ff"`
}

type FixtureA struct {
	XMLName       struct{}      `xml:"BixResponse"`
	FromNamespace string        `xml:"fromNamespace"`
	Request       string        `xml:"request"`
	MOID          uint64        `xml:"moid"`
	MessageID     uint64        `xml:"mid"`
	Data          FixtureA_Data `xml:"Data"`
}

type FixtureA_Data struct {
	XMLName struct{} `xml:"Data"`
	Auth    bool     `xml:"auth"`
}

type FixtureB struct {
	XMLName   struct{}           `xml:"FixtureB"`
	StringMap FixtureB_StringMap `xml:"StringMap"`
}

type FixtureB_StringMap map[string]string

func (s FixtureB_StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{start}
	for key, value := range s {
		t := xml.StartElement{Name: xml.Name{"", key}}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{t.Name})
	}
	tokens = append(tokens, xml.EndElement{start.Name})
	for _, t := range tokens {
		if err := e.EncodeToken(t); err != nil {
			return err
		}
	}
	return e.Flush()
}

func TestEncodeFixture1(t *testing.T) {
	assert := assert.New(t)

	// Read XML fixture #1
	fixture := "testdata/test-systemlib-1.xml"
	fmt.Printf("Loading fixture %s\n", fixture)
	xmlBytes, err := ioutil.ReadFile(fixture)
	assert.NoError(err)

	// Unmarshal XML file into Fixture1 structure
	fixture1 := Fixture1{}
	err = xml.Unmarshal(xmlBytes, &fixture1)
	assert.NoError(err)

	// Sanity check - ensure fixture1 contains expected values loaded from file, before starting test
	assert.Equal("VirtualMachines", fixture1.ToNamespace)
	assert.Equal("Testing", fixture1.Request)
	assert.Equal("6", fixture1.MOID)
	assert.Equal("1", fixture1.MID)

	// Encode structure as binary xml
	file, _ := ioutil.TempFile("", "binaryxmlEncoderTest1")
	writer := bufio.NewWriter(file)
	fmt.Printf("Writing binary encoded xml file %s\n", file.Name())
	err = binaryxml.Encode(fixture1, writer)
	assert.NoError(err)
	writer.Flush()

	// Unmarshal binary XML file into 2nd Fixture1 structure
	binaryXmlBytes, err := ioutil.ReadFile(file.Name())
	assert.NoError(err)
	xmlString, err := binaryxml.ToXML(binaryXmlBytes)
	assert.NoError(err)
	secondFixture1 := Fixture1{}
	err = xml.Unmarshal([]byte(xmlString), &secondFixture1)
	assert.NoError(err)

	// Perform test
	assert.Equal(fixture1, secondFixture1)
}

func TestEncodeFixture2(t *testing.T) {
	assert := assert.New(t)

	// Read XML fixture #2
	fixture := "testdata/test-systemlib-2.xml"
	fmt.Printf("Loading fixture %s\n", fixture)
	xmlBytes, err := ioutil.ReadFile(fixture)
	assert.NoError(err)

	// Unmarshal XML file into Fixture2 structure
	fixture2 := Fixture2{}
	err = xml.Unmarshal(xmlBytes, &fixture2)
	assert.NoError(err)

	// Sanity check - ensure fixture2 contains expected values loaded from file, before starting test
	assert.Equal("SubscriptionProvider", fixture2.ToNamespace)
	assert.Equal("Subscribe", fixture2.Request)
	assert.Equal("2", fixture2.MOID)
	assert.Equal("2", fixture2.MID)
	assert.Equal("3", fixture2.Data.MOID)
	assert.Equal("Common_CPU", fixture2.Data.Query.Namespace)
	assert.Equal("_Total", fixture2.Data.Query.Instance)
	assert.Equal("*", fixture2.Data.Query.Key)
	assert.Equal("IdleTime", fixture2.Data.Query.Field)
	assert.Equal(uint32(60), fixture2.Data.Query.Interval)
	assert.Equal("", fixture2.Data.Query.UserKey)
	assert.Equal("SB", fixture2.Data.Query.Tablespace)

	// Encode structure as binary xml
	file, _ := ioutil.TempFile("", "binaryxmlEncoderTest2")
	writer := bufio.NewWriter(file)
	fmt.Printf("Writing binary encoded xml file %s\n", file.Name())
	err = binaryxml.Encode(fixture2, writer)
	assert.NoError(err)
	writer.Flush()

	// Unmarshal binary XML file into 2nd Fixture2 structure
	binaryXmlBytes, err := ioutil.ReadFile(file.Name())
	assert.NoError(err)
	xmlString, err := binaryxml.ToXML(binaryXmlBytes)
	assert.NoError(err)
	secondFixture2 := Fixture2{}
	err = xml.Unmarshal([]byte(xmlString), &secondFixture2)
	assert.NoError(err)

	// Perform test
	assert.Equal(fixture2, secondFixture2)
}

func TestEncodeFixture3(t *testing.T) {
	assert := assert.New(t)

	// Read XML fixture #3
	fixture := "testdata/test-systemlib-3.xml"
	fmt.Printf("Loading fixture %s\n", fixture)
	xmlBytes, err := ioutil.ReadFile(fixture)
	assert.NoError(err)

	// Unmarshal XML file into Fixture2 structure
	fixture3 := Fixture3{}
	err = xml.Unmarshal(xmlBytes, &fixture3)
	assert.NoError(err)

	// Sanity check - ensure fixture3 contains expected values loaded from file, before starting test
	assert.Equal("Waveforms", fixture3.Namespace.Name)
	assert.Equal("Waveforms", fixture3.Namespace.Display)
	assert.Equal("1", fixture3.Namespace.Version)
	assert.Equal("schema", fixture3.Namespace.Schema.Name)
	assert.Equal("default", fixture3.Namespace.Schema.Instance.Name)
	assert.Equal("default", fixture3.Namespace.Schema.Instance.Display)
	field := Fixture3_Field{Name: "sinewave", Type: "schema_uint32", Display: "sinewave", UserData: "Custom field"}
	assert.Contains(fixture3.Namespace.Schema.Fields, field)
	field = Fixture3_Field{Name: "random", Type: "schema_uint32", Display: "random", UserData: "Custom field"}
	assert.Contains(fixture3.Namespace.Schema.Fields, field)

	// Encode structure as binary xml
	file, _ := ioutil.TempFile("", "binaryxmlEncoderTest3")
	writer := bufio.NewWriter(file)
	fmt.Printf("Writing binary encoded xml file %s\n", file.Name())
	err = binaryxml.Encode(fixture3, writer)
	assert.NoError(err)
	writer.Flush()

	// Unmarshal binary XML file into 2nd Fixture3 structure
	binaryXmlBytes, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("Failed opening generated binary xml file %s", file.Name)
	}
	xmlString, err := binaryxml.ToXML(binaryXmlBytes)
	assert.NoError(err)
	secondFixture3 := Fixture3{}
	err = xml.Unmarshal([]byte(xmlString), &secondFixture3)
	assert.NoError(err)

	// Perform test
	assert.Equal(fixture3, secondFixture3)
}

func TestEncodeFixture4(t *testing.T) {
	assert := assert.New(t)

	// Read XML fixture #4
	fixture := "testdata/test-systemlib-4.xml"
	fmt.Printf("Loading fixture %s\n", fixture)
	xmlBytes, err := ioutil.ReadFile(fixture)
	assert.NoError(err)

	// Unmarshal XML file into Fixture4 structure
	fixture4 := Fixture4{}
	err = xml.Unmarshal(xmlBytes, &fixture4)
	assert.NoError(err)

	// Sanity check - ensure fixture4 contains expected values loaded from file, before starting test
	assert.Equal(int8(-128), fixture4.Int8Min)
	assert.Equal(int8(127), fixture4.Int8Max)
	assert.Equal(int16(-32768), fixture4.Int16Min)
	assert.Equal(int16(32767), fixture4.Int16Max)
	assert.Equal(int32(-2147483648), fixture4.Int32Min)
	assert.Equal(int32(2147483647), fixture4.Int32Max)
	assert.Equal(int64(-9223372036854775808), fixture4.Int64Min)
	assert.Equal(int64(9223372036854775807), fixture4.Int64Max)
	assert.Equal(uint8(255), fixture4.Uint8Max)
	assert.Equal(uint16(65535), fixture4.Uint16Max)
	assert.Equal(uint32(4294967295), fixture4.Uint32Max)
	assert.Equal(uint64(18446744073709551615), fixture4.Uint64Max)

	// Encode structure as binary xml
	file, _ := ioutil.TempFile("", "binaryxmlEncoderTest4")
	writer := bufio.NewWriter(file)
	fmt.Printf("Writing binary encoded xml file %s\n", file.Name())
	err = binaryxml.Encode(fixture4, writer)
	assert.NoError(err)
	writer.Flush()

	// Unmarshal binary XML file into 2nd Fixture4 structure
	binaryXmlBytes, err := ioutil.ReadFile(file.Name())
	assert.NoError(err)
	xmlString, err := binaryxml.ToXML(binaryXmlBytes)
	assert.NoError(err)
	secondFixture4 := Fixture4{}
	err = xml.Unmarshal([]byte(xmlString), &secondFixture4)
	assert.NoError(err)

	// Perform test
	assert.Equal(fixture4, secondFixture4)
}

func TestEncodeFixture5(t *testing.T) {
	assert := assert.New(t)

	// Read XML fixture #5
	fixture := "testdata/test-systemlib-5.xml"
	fmt.Printf("Loading fixture %s\n", fixture)
	xmlBytes, err := ioutil.ReadFile(fixture)
	assert.NoError(err)

	// Unmarshal XML file into Fixture5 structure
	fixture5 := Fixture5{}
	err = xml.Unmarshal(xmlBytes, &fixture5)
	assert.NoError(err)

	// Sanity check - ensure fixture5 contains expected values loaded from file, before starting test
	assert.Equal(float32(0), fixture5.Float32_0)
	assert.Equal(float32(3.14), fixture5.Float32_Pi)
	assert.Equal(float32(-3.14), fixture5.Float32_NegativePi)

	// Encode structure as binary xml
	file, _ := ioutil.TempFile("", "binaryxmlEncoderTest5")
	writer := bufio.NewWriter(file)
	fmt.Printf("Writing binary encoded xml file %s\n", file.Name())
	err = binaryxml.Encode(fixture5, writer)
	assert.NoError(err)
	writer.Flush()

	// Unmarshal binary XML file into 2nd Fixture5 structure
	binaryXmlBytes, err := ioutil.ReadFile(file.Name())
	assert.NoError(err)
	xmlString, err := binaryxml.ToXML(binaryXmlBytes)
	assert.NoError(err)
	secondFixture5 := Fixture5{}
	err = xml.Unmarshal([]byte(xmlString), &secondFixture5)
	assert.NoError(err)

	// Perform test
	assert.Equal(fixture5, secondFixture5)
}

func TestEncodeFixture6(t *testing.T) {
	assert := assert.New(t)

	// Statically prepare Fixture6 structure, and sanity check
	fixture6 := Fixture6{Binary_007f80ff: []byte{0x00, 0x7f, 0x80, 0xff}}

	// Encode structure as binary xml
	file, _ := ioutil.TempFile("", "binaryxmlEncoderTest6")
	writer := bufio.NewWriter(file)
	fmt.Printf("Writing binary encoded xml file %s\n", file.Name())
	err := binaryxml.Encode(fixture6, writer)
	assert.NoError(err)
	writer.Flush()

	// Load binary xml files and compare (they're single-field so exact matching can be used)
	expectedBinaryXmlBytes, err := ioutil.ReadFile("testdata/test-systemlib-6.binaryxml")
	assert.NoError(err)
	actualBinaryXmlBytes, err := ioutil.ReadFile(file.Name())
	assert.NoError(err)

	// Perform test
	assert.Equal(expectedBinaryXmlBytes, actualBinaryXmlBytes)
}

func TestEncodeFixtureA(t *testing.T) {
	assert := assert.New(t)

	// Statically prepare FixtureA structure, and sanity check
	fixtureA := FixtureA{FromNamespace: "_internal", Request: "_GETAUTH", MOID: 18446744073709551615, MessageID: 27}
	fixtureA.Data.Auth = false

	// Encode structure as binary xml
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err := binaryxml.Encode(fixtureA, writer)
	assert.NoError(err)
	writer.Flush()
	actualBinaryXML := b.Bytes()

	// Marshal fixtureA to XML
	expectedXMLBytes, err := xml.Marshal(fixtureA)
	assert.NoError(err)
	expectedXML := string(expectedXMLBytes)

	// Convert binary xml to XML
	actualXMLBytes, err := binaryxml.ToXML([]byte(actualBinaryXML))
	assert.NoError(err)
	actualXML := string(actualXMLBytes)

	// Minify expected vs actual xml, to minimize inconsequential differences, then compare
	minifier := minify.New()
	minifier.AddFunc("text/xml", xmlminifier.Minify)
	minifiedExpectedXML, err := minifier.String("text/xml", expectedXML)
	assert.NoError(err)
	minifiedActualXML, err := minifier.String("text/xml", actualXML)
	assert.NoError(err)
	assert.Equal(minifiedExpectedXML, minifiedActualXML)
}

// ----------------------------------------------------------------------------
// TestEncodeFixtureB
// ----------------------------------------------------------------------------

func TestEncodeFixtureB(t *testing.T) {
	assert := assert.New(t)

	// Prepare test data
	var fixture FixtureB
	fixture.StringMap = make(FixtureB_StringMap)
	fixture.StringMap["abc"] = "123"
	const expected = "<FixtureB><StringMap><abc>123</abc></StringMap></FixtureB>"

	// Sanity check - struct to XML via Go Marshaller
	actualBytes, err := xml.Marshal(fixture)
	assert.NoError(err)
	actual := string(actualBytes)
	assert.Equal(expected, actual)

	// Test -- struct to BinaryXML to XML
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	assert.NoError(binaryxml.Encode(fixture, writer))
	writer.Flush()
	binaryXML := buffer.Bytes()
	actual, err = binaryxml.ToXML(binaryXML)
	assert.NoError(err)
	assert.Equal(expected, actual)
}
