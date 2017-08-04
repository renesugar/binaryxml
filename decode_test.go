package binaryxml_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/BixData/binaryxml"
	"github.com/stretchr/testify/assert"
)

func TestDecodeFixture1(t *testing.T) {
	assert := assert.New(t)

	// Read Binary XML fixture #1
	fixture := "testdata/test-systemlib-1.binaryxml"
	fmt.Printf("Loading fixture %s\n", fixture)
	binaryXML, err := ioutil.ReadFile(fixture)
	assert.NoError(err)

	// Decode XML file into Fixture1 structure
	fixture1 := Fixture1{}
	assert.NoError(binaryxml.Decode(binaryXML, &fixture1))

	// Ensure fixture1 contains expected values loaded from file, before starting test
	assert.Equal("VirtualMachines", fixture1.ToNamespace)
	assert.Equal("Testing", fixture1.Request)
	assert.Equal("6", fixture1.MOID)
	assert.Equal("1", fixture1.MID)
}

func TestDecodeFixture2(t *testing.T) {
	assert := assert.New(t)

	// Read Binary XML fixture #2
	fixture := "testdata/test-systemlib-2.binaryxml"
	fmt.Printf("Loading fixture %s\n", fixture)
	binaryXML, err := ioutil.ReadFile(fixture)
	assert.NoError(err)

	// Decode XML file into Fixture2 structure
	fixture2 := Fixture2{}
	assert.NoError(binaryxml.Decode(binaryXML, &fixture2))

	// Ensure fixture2 contains expected values loaded from file, before starting test
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
}

func TestDecodeFixture3(t *testing.T) {
	assert := assert.New(t)

	// Read Binary XML fixture #3
	fixture := "testdata/test-systemlib-3.binaryxml"
	fmt.Printf("Loading fixture %s\n", fixture)
	binaryXML, err := ioutil.ReadFile(fixture)
	assert.NoError(err)

	// Decode Binary XML into Fixture2 structure
	fixture3 := Fixture3{}
	assert.NoError(binaryxml.Decode(binaryXML, &fixture3))

	// Ensure fixture3 contains expected values loaded from file, before starting test
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
}

func TestDecodeFixture4(t *testing.T) {
	assert := assert.New(t)

	// Read Binary XML fixture #4
	fixture := "testdata/test-systemlib-4.binaryxml"
	fmt.Printf("Loading fixture %s\n", fixture)
	binaryXML, err := ioutil.ReadFile(fixture)
	assert.NoError(err)

	// Decode Binary XML into Fixture4 structure
	fixture4 := Fixture4{}
	assert.NoError(binaryxml.Decode(binaryXML, &fixture4))

	// Ensure fixture4 contains expected values loaded from file, before starting test
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
}

func TestDecodeFixture5(t *testing.T) {
	assert := assert.New(t)

	// Read Binary XML fixture #5
	fixture := "testdata/test-systemlib-5.binaryxml"
	fmt.Printf("Loading fixture %s\n", fixture)
	binaryXML, err := ioutil.ReadFile(fixture)
	assert.NoError(err)

	// Decode Binary XML into Fixture5 structure
	fixture5 := Fixture5{}
	assert.NoError(binaryxml.Decode(binaryXML, &fixture5))

	// Ensure fixture5 contains expected values loaded from file, before starting test
	assert.Equal(float32(0), fixture5.Float32_0)
	assert.Equal(float32(3.14), fixture5.Float32_Pi)
	assert.Equal(float32(-3.14), fixture5.Float32_NegativePi)
}
