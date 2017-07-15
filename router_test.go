package binaryxml

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

// ----------------------------------------------------------------------------

func TestRegisterRoutes(t *testing.T) {
	assert := assert.New(t)

	router := NewRouter()
	assert.NotNil(router)

	router.Add("/BixRequest[toNamespace='VirtualMachines'][request='Testing']", func(*Context) error {
		return nil
	})
	binaryXml, err := ioutil.ReadFile("testdata/test-systemlib-1.binaryxml")
	assert.NoError(err)
	request, err := NewRequest(binaryXml)
	ctx := NewContext(request)
	handler := router.findHandler(ctx)
	assert.NotNil(handler)
}

// ----------------------------------------------------------------------------

func TestRouteFixture1(t *testing.T) {
	assert := assert.New(t)

	router := NewRouter()
	assert.NotNil(router)

	handlerCalled := false
	router.Add("/BixRequest[toNamespace='VirtualMachines'][request='Testing']", func(ctx *Context) error {
		ctx.Response.BinaryXML = []byte{byte(tablebegin), byte(tableend), byte(serialbegin), byte(serialend)}
		handlerCalled = true
		return nil
	})
	binaryXml, err := ioutil.ReadFile("testdata/test-systemlib-1.binaryxml")
	assert.NoError(err)
	request, err := NewRequest(binaryXml)
	assert.Equal(uint64(6), request.MOID())
	assert.Equal(uint64(1), request.MID())
	assert.Equal("Testing", request.Request())
	assert.Equal("VirtualMachines", request.ToNamespace())
	ctx := NewContext(request)

	router.Handle(ctx)
	assert.Equal(true, handlerCalled)
	assert.Equal(4, len(ctx.Response.BinaryXML))
}

// ----------------------------------------------------------------------------

func TestSetResponseError(t *testing.T) {
	assert := assert.New(t)
	binaryXml, err := ioutil.ReadFile("testdata/test-systemlib-1.binaryxml")
	assert.NoError(err)
	request, err := NewRequest(binaryXml)
	assert.NoError(err)
	ctx := NewContext(request)
	ctx.Response.Error(request, "foo")
	xml, err := ToXML(ctx.Response.BinaryXML)
	expected := "<?xml version=\"1.0\"?>\n<BixError><fromNamespace>VirtualMachines</fromNamespace><request>Testing</request><moid>6</moid><mid>1</mid><error>foo</error></BixError>"
	assert.Equal(expected, xml)
}
