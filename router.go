package binaryxml

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/docktermj/go-logger/logger"
	"github.com/jnewmoyer/xmlpath"
)

// ----------------------------------------------------------------------------
// Router request
// ----------------------------------------------------------------------------

type Request struct {
	RemoteAddr  string
	XML         string
	BinaryXML   []byte
	Param       uint8
	XMLPathNode *xmlpath.Node
}

func NewRequest(binaryXml []byte) (*Request, error) {
	// Populate XML field
	xml, err := ToXML(binaryXml)
	if err != nil {
		return nil, err
	}

	// Parse XML for future xpath queries
	xmlPathNode, err := xmlpath.Parse(strings.NewReader(xml))
	if err != nil {
		return nil, err
	}

	request := Request{BinaryXML: binaryXml, XML: xml, XMLPathNode: xmlPathNode}
	return &request, nil
}

func (request *Request) MID() uint64 {
	path := xmlpath.MustCompile("/BixRequest/mid")
	if value, ok := path.String(request.XMLPathNode); ok {
		if mid, err := strconv.ParseUint(value, 10, 64); err == nil {
			return mid
		}
	}
	return 0
}

func (request *Request) MOID() uint64 {
	path := xmlpath.MustCompile("/BixRequest/moid")
	if value, ok := path.String(request.XMLPathNode); ok {
		if moid, err := strconv.ParseUint(value, 10, 64); err == nil {
			return moid
		}
	}
	return 0
}

func (request *Request) Name() string {
	path := xmlpath.MustCompile("/")
	iter := path.Iter(request.XMLPathNode)
	for iter.Next() {
		node := iter.Node()
		return node.Name().Local
	}
	return ""
}

func (request *Request) Namespace() string {
	path := xmlpath.MustCompile("/BixRequest/toNamespace")
	if value, ok := path.String(request.XMLPathNode); ok {
		return value
	}
	return ""
}

func (request *Request) Request() string {
	path := xmlpath.MustCompile("/BixRequest/request")
	if value, ok := path.String(request.XMLPathNode); ok {
		return value
	}
	return ""
}

// ----------------------------------------------------------------------------
// Router response
// ----------------------------------------------------------------------------

type Response struct {
	BinaryXML []byte
	Param     uint8
}

// ----------------------------------------------------------------------------
// Router request context
// ----------------------------------------------------------------------------

type SendMoreFunc func(*Context) error

type Context struct {
	Request      *Request
	Response     *Response
	SendMoreFunc SendMoreFunc
}

func NewContext(request *Request) *Context {
	response := &Response{}
	return &Context{Request: request, Response: response}
}

func (ctx *Context) Respond(v interface{}) error {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	if err := Encode(v, writer); err != nil {
		return err
	}
	writer.Flush()
	ctx.Response.BinaryXML = b.Bytes()
	return nil
}

func (ctx *Context) RespondMore(v interface{}) error {
	if err := ctx.Respond(v); err != nil {
		return err
	}
	return ctx.sendMore()
}

func (ctx *Context) RespondError(message string) error {
	bixError := BixError{FromNamespace: ctx.Request.Namespace(), Request: ctx.Request.Request(), MOID: ctx.Request.MOID(), MID: ctx.Request.MID(), Error: message}
	return ctx.Respond(bixError)
}

func (ctx *Context) sendMore() error {
	if ctx.SendMoreFunc == nil {
		return errors.New("Cannot RespondMore without a Context SendMoreFunc set")
	}
	return ctx.SendMoreFunc(ctx)
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

// ----------------------------------------------------------------------------

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
		path := xmlpath.MustCompile(xpath)
		if _, ok := path.String(ctx.Request.XMLPathNode); ok {
			handler := router.registry[xpath]
			return handler
		}
	}
	return router.defaultHandler
}

func (router *routerImpl) Handle(ctx *Context) error {
	handler := router.findHandler(ctx)
	topic := fmt.Sprintf("%s %s::%s", ctx.Request.Name(), ctx.Request.Namespace(), ctx.Request.Request())
	if handler == nil {
		logger.Warnf("No handler for %s", topic)
		return nil
	}
	err := handler(ctx)
	return err
}
