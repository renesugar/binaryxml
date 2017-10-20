package binaryxml_router

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/BixData/binaryxml"
	"github.com/docktermj/go-logger/logger"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds | log.LUTC)
	logger.SetLevel(logger.LevelInfo)

	os.Exit(m.Run())
}

// ----------------------------------------------------------------------------

func TestRegisterRoutes(t *testing.T) {
	assert := assert.New(t)

	router := NewRouter()
	assert.NotNil(router)

	router.Add("/BixRequest[toNamespace='VirtualMachines' and request='Testing']", func(*Context) error {
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

	testingHandlerCalled := false
	otherHandler1Called := false
	otherHandler2Called := false
	router.Add("/BixRequest[toNamespace='VirtualMachines' and request='1']", func(ctx *Context) error {
		otherHandler1Called = true
		return nil
	})
	router.Add("/BixRequest[toNamespace='VirtualMachines' and request='Testing']", func(ctx *Context) error {
		testingHandlerCalled = true
		return nil
	})
	router.Add("/BixRequest[toNamespace='VirtualMachines' and request='Z']", func(ctx *Context) error {
		otherHandler2Called = true
		return nil
	})
	binaryXml, err := ioutil.ReadFile("testdata/test-systemlib-1.binaryxml")
	assert.NoError(err)
	request, err := NewRequest(binaryXml)
	assert.Equal("BixRequest", request.Name())
	assert.Equal(uint64(6), request.MOID())
	assert.Equal(uint64(1), request.MID())
	assert.Equal("VirtualMachines", request.Namespace())
	assert.Equal("Testing", request.Request())
	ctx := NewContext(request)

	router.Handle(ctx)
	assert.Equal(true, testingHandlerCalled)
	assert.Equal(false, otherHandler1Called)
	assert.Equal(false, otherHandler2Called)
}

// ----------------------------------------------------------------------------

func TestSetResponseError(t *testing.T) {
	assert := assert.New(t)
	binaryXml, err := ioutil.ReadFile("testdata/test-systemlib-1.binaryxml")
	assert.NoError(err)
	request, err := NewRequest(binaryXml)
	assert.NoError(err)
	ctx := NewContext(request)
	ctx.RespondError("foo")
	xml, err := binaryxml.ToXML(ctx.Response.BinaryXML)
	expected := "<BixError><fromNamespace>VirtualMachines</fromNamespace><request>Testing</request><moid>6</moid><mid>1</mid><error>foo</error></BixError>"
	assert.Equal(expected, xml)
}

// ----------------------------------------------------------------------------
// Test response continuation capability.
//
//   Recv: BixRequest
//   Send: BixResponse - RespondMore() // intermediate "more" response
//   Send: BixResponse - RespondMore() // intermediate "more" response
//   Send: BixResponse - Respond       // final response
// ----------------------------------------------------------------------------

func TestRespondMore(t *testing.T) {
	assert := assert.New(t)
	type bixResponse struct {
		XMLName struct{} `xml:"BixResponse"`
		Data    string   `xml:"Data"`
	}
	binaryXml, err := ioutil.ReadFile("testdata/test-systemlib-1.binaryxml")
	assert.NoError(err)
	request, err := NewRequest(binaryXml)
	assert.NoError(err)
	ctx := NewContext(request)
	assert.Error(ctx.RespondMore(bixResponse{Data: "partial"}), "Expected error if no 'send more' handler has been set")
	sendMoreFuncCalled := false
	ctx.SendMoreFunc = func(ctx *Context) error {
		sendMoreFuncCalled = true
		return nil
	}
	assert.NoError(ctx.RespondMore(bixResponse{Data: "partial"}))
	ctx.Respond(bixResponse{Data: "done"})
	xml, err := binaryxml.ToXML(ctx.Response.BinaryXML)
	expected := "<BixResponse><Data>done</Data></BixResponse>"
	assert.Equal(expected, xml)
}
