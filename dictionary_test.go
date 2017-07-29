package binaryxml

import (
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/cevaris/ordered_map"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Test with Fixture1
// ----------------------------------------------------------------------------

type Fixture1 struct {
	XMLName     struct{} `xml:"BixRequest"`
	Request     string   `xml:"request"`
	ToNamespace string   `xml:"toNamespace"`
	MOID        string   `xml:"moid"`
	MID         string   `xml:"mid"`
}

func TestGenerateDictionaryForFixture1(t *testing.T) {
	assert := assert.New(t)
	fixture := Fixture1{}
	dictionary := ordered_map.NewOrderedMap()
	assert.NoError(generateElementNameDictionaryForValue(reflect.ValueOf(fixture), dictionary))
	assert.Equal(5, dictionary.Len())
	assert.NotEmpty(dictionary.Get("BixRequest"))
	assert.NotEmpty(dictionary.Get("request"))
	assert.NotEmpty(dictionary.Get("toNamespace"))
	assert.NotEmpty(dictionary.Get("moid"))
	assert.NotEmpty(dictionary.Get("mid"))
	assert.Empty(dictionary.Get("bogus"))
}

// ----------------------------------------------------------------------------
// Test with FixtureB
// ----------------------------------------------------------------------------

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

func TestGenerateDictionaryForFixtureB(t *testing.T) {
	assert := assert.New(t)
	fixture := FixtureB{}
	fixture.StringMap = make(FixtureB_StringMap)
	fixture.StringMap["abc"] = "123"

	dictionary := ordered_map.NewOrderedMap()
	assert.NoError(generateElementNameDictionaryForValue(reflect.ValueOf(fixture), dictionary))
	assert.NotEmpty(dictionary.Get("StringMap"))
	assert.NotEmpty(dictionary.Get("abc"))
	assert.Empty(dictionary.Get("123"))
}
