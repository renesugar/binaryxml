package binaryxml

import (
	"github.com/antchfx/xquery/xml"
	"strings"
)


//=======================================================================
// Router request
//=======================================================================

type Request struct {
	RemoteAddr string
	XML string
	BinaryXML []byte
	XMLQueryNode *xmlquery.Node
}


func NewRequest(binaryXml []byte) (*Request, error) {
	// Populate XML field
	xml, err := ToXML(binaryXml)
	if err != nil {return nil, err}
	
	// Parse XML for future xpath queries
	xmlQueryNode, err := xmlquery.Parse(strings.NewReader(xml))
	if err != nil {return nil, err}
	
	request := Request{BinaryXML:binaryXml, XML:xml, XMLQueryNode:xmlQueryNode}
	return &request, nil
}


//=======================================================================
// Router response
//=======================================================================

type Response struct {
	BinaryXML []byte
}


//=======================================================================
// Router request context
//=======================================================================

type Context struct {
	Request *Request
	Response *Response
}


func NewContext(request *Request) *Context {
	response := &Response{}
	return &Context{Request:request, Response:response}
}


//=======================================================================
// Router route handler
//=======================================================================

type HandlerFunc func(*Context) error


//=======================================================================
// Router
//=======================================================================

type Router interface {
	// Register a handler for a given xpath
	Add(xpath string, handler HandlerFunc)
	
	// Register a default handler for use when no others are registered for a given xpath
	Default(handler HandlerFunc)
	
	// Find a handler function to match the given request
	findHandler(ctx *Context) HandlerFunc
	
	Handle(ctx *Context)
}


type routerImpl struct {
	registry map[string]HandlerFunc
	defaultHandler HandlerFunc
}


func NewRouter() *routerImpl {
  return &routerImpl{registry:make(map[string]HandlerFunc)}
}


func (router *routerImpl) Add(xpath string, handler HandlerFunc) {
	router.registry[xpath] = handler
}


func (router *routerImpl) Default(handler HandlerFunc) {
	router.defaultHandler = handler
}


func (router *routerImpl) findHandler(ctx *Context) HandlerFunc {
	for xpath := range router.registry {
    	if node := xmlquery.FindOne(ctx.Request.XMLQueryNode, xpath); node != nil {
			handler := router.registry[xpath]
			return handler
    	}
	}
	return router.defaultHandler
}


func (router *routerImpl) Handle(ctx *Context) error {
	handler := router.findHandler(ctx)
	if handler == nil {return nil}
	err := handler(ctx)
	return err
}
