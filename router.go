package binaryxml

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	"github.com/antchfx/xquery/xml"
)

// ----------------------------------------------------------------------------
// Router request
// ----------------------------------------------------------------------------

type Request struct {
	RemoteAddr   string
	XML          string
	BinaryXML    []byte
	XMLQueryNode *xmlquery.Node
}

func NewRequest(binaryXml []byte) (*Request, error) {
	// Populate XML field
	xml, err := ToXML(binaryXml)
	if err != nil {
		return nil, err
	}

	// Parse XML for future xpath queries
	xmlQueryNode, err := xmlquery.Parse(strings.NewReader(xml))
	if err != nil {
		return nil, err
	}

	request := Request{BinaryXML: binaryXml, XML: xml, XMLQueryNode: xmlQueryNode}
	return &request, nil
}

func (request *Request) MID() uint64 {
	node := xmlquery.FindOne(request.XMLQueryNode, "/BixRequest/mid")
	if node == nil {
		return 0
	}
	mid, err := strconv.ParseUint(node.InnerText(), 10, 64)
	if err != nil {
		return 0
	}
	return mid
}

func (request *Request) MOID() uint64 {
	node := xmlquery.FindOne(request.XMLQueryNode, "/BixRequest/moid")
	if node == nil {
		return 0
	}
	moid, err := strconv.ParseUint(node.InnerText(), 10, 64)
	if err != nil {
		return 0
	}
	return moid
}

func (request *Request) Request() string {
	node := xmlquery.FindOne(request.XMLQueryNode, "/BixRequest/request")
	if node == nil {
		return ""
	}
	return node.InnerText()
}

func (request *Request) ToNamespace() string {
	node := xmlquery.FindOne(request.XMLQueryNode, "/BixRequest/toNamespace")
	if node == nil {
		return ""
	}
	return node.InnerText()
}

// ----------------------------------------------------------------------------
// Router response
// ----------------------------------------------------------------------------

type Response struct {
	BinaryXML []byte
}

func (response *Response) Encode(v interface{}) error {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	if err := Encode(v, writer); err != nil {
		return err
	}
	writer.Flush()
	response.BinaryXML = b.Bytes()
	return nil
}

func (response *Response) Error(request *Request, message string) error {
	bixError := BixError{FromNamespace: request.ToNamespace(), Request: request.Request(), MOID: request.MOID(), MID: request.MID(), Error: message}
	return response.Encode(bixError)
}

// ----------------------------------------------------------------------------
// Router request context
// ----------------------------------------------------------------------------

type Context struct {
	Request  *Request
	Response *Response
}

func NewContext(request *Request) *Context {
	response := &Response{}
	return &Context{Request: request, Response: response}
}

// ----------------------------------------------------------------------------
// Router route handler
// ----------------------------------------------------------------------------

type HandlerFunc func(*Context) error

// ----------------------------------------------------------------------------
// Router
// ----------------------------------------------------------------------------

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
	registry       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

func NewRouter() *routerImpl {
	return &routerImpl{registry: make(map[string]HandlerFunc)}
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
	if handler == nil {
		return nil
	}
	err := handler(ctx)
	return err
}
