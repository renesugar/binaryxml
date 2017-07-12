package binaryxml

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)


func TestRegisterRoutes(t *testing.T) {
	assert := assert.New(t)
	
	router := NewRouter()
	assert.NotNil(router)
	
	router.Add("/BixRequest[toNamespace='VirtualMachines'][request='Testing']", func(*Context) error {
		return nil
	})
	binaryXml, err := ioutil.ReadFile("testdata/test-systemlib-1.binaryxml")
	assert.Nil(err)
	request, err := NewRequest(binaryXml)
	ctx := NewContext(request)
	handler := router.findHandler(ctx)
	assert.NotNil(handler)
}


func TestRouteFixture1(t *testing.T) {
	assert := assert.New(t)
	
	router := NewRouter()
	assert.NotNil(router)
	
	handlerCalled := false
	router.Add("/BixRequest[toNamespace='VirtualMachines'][request='Testing']", func(*Context) error {
		handlerCalled = true
		return nil
	})
	binaryXml, err := ioutil.ReadFile("testdata/test-systemlib-1.binaryxml")
	assert.Nil(err)
	request, err := NewRequest(binaryXml)
	ctx := NewContext(request)
	
	router.Handle(ctx)
	assert.Equal(true, handlerCalled)
}
