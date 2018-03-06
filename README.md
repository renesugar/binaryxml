# Binary XML Library

This GoLang library provides encoding and decoding support for NuBix binary xml.

## Table of Contents

* [Convert Binary XML to XML](#convert-binary-xml-to-xml)
* [Encode a Struct](#encode-a-struct)
* [Decode a Struct](#decode-a-struct)
* [Routing](#routing)
  * [Routing Requests](#routing-requests)
* [Testing](#testing)

## Convert Binary XML to XML

```go
import (
	"github.com/BixData/binaryxml"
	"io/ioutil"
)

binaryXml, _ := ioutil.ReadFile("mydata.binaryxml")
xml, err := binaryxml.ToXML(binaryXml)
```

## Encode a Struct

The following code converts a struct to Binary XML.

```go
import (
	"bufio"
	"encoding/xml"
	"github.com/BixData/binaryxml"
	"io/ioutil"
)

type Person struct {
	XMLName struct{} `xml:"Person"`
	First   string   `xml:"firstName"`
	Last    string   `xml:"lastName"`
	Age     uint8    `xml:"age"`
}

person := Person{First:"John", Last:"Doe", Age:49}

file, _ := ioutil.TempFile("", "myBinaryXmlFile")
writer := bufio.NewWriter(file)
err := binaryxml.Encode(person, writer);
writer.Flush()
```

## Decode a Struct

Hydrating a struct with decoded Binary XML is currently accomplished through intermediate use of XML, which might experience some loss of data fidelity due to sub-optimal datatype transfer. This feature is ripe for future improvement.

```go
person := Person{}
err := binaryxml.Decode(binaryXml, &person)
```

## Routing

The `router` sub-package provides a network reactor that assigns incoming messages to handlers according to XPath expressions  designed to be matched against BixRequest fields. This is meant to provide a more modern alternative to the Bix `MessageObject` peering interface. This package is made separate so that it can be ignored, if a pure Bix `MessageObject` reactor will be used instead.

### Routing Requests

```go
router := binaryxml.NewRouter()
router.Add("/BixRequest[toNamespace='SubscriptionManager'][request='Subscribe']", handleSubscribeRequest)
router.Add("/BixRequest[toNamespace='_internal'][request='_GETAUTH']", handleInternalGetAuthRequest)

func handleInternalGetAuthRequest(ctx *Context) error {
	// Prepare response object
	type bixResponse struct {
		XMLName       struct{}        `xml:"BixResponse"`
		FromNamespace string          `xml:"fromNamespace"`
		Request       string          `xml:"request"`
		MOID          uint64          `xml:"moid"`
		MID           uint64          `xml:"mid"`
		Data          bixResponseData `xml:"Data"`
	}
	type bixResponseData struct {
		XMLName struct{} `xml:"Data"`
		Auth    bool     `xml:"auth"`
	}
	bixRes := bixResponse{FromNamespace:"_internal", Request:'_GETAUTH'}
	bixRes.Data.Auth = false

	// Serialize response object to binaryxml
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err := binaryxml.Encode(bixRes, writer)
	writer.Flush()
	binaryXML := b.Bytes()

	// Send response
	ctx.Response.BinaryXML = binaryXML
	return nil
})

listener, err := net.Listen("tcp", 17070)
conn, err := listener.Accept()
reader := bufio.NewReader(conn)
buffer := make([]byte, 100000)
for {
	readBytes, err := reader.Read(buffer)
	binaryXml := buf[:readBytes]
	request, err := binaryxml.NewRequest(binaryXml)
	ctx := binaryxml.NewContext(request)
	router.Handle(ctx)
	conn.Write(ctx.Response.BinaryXML)
}
```

## Testing

Setup a workspace:

```sh
$ mkdir -p workspaces/go/src/github.com/BixData
$ cd workspaces/go
$ export GOPATH=`pwd`
$ cd src/github.com/BixData
$ git clone <this repo>
$ cd <this repo>
$ make dependencies
```

And then test with:

```sh
$ go test ./...
ok  	github.com/BixData/binaryxml	0.038s
ok  	github.com/BixData/binaryxml/router	0.033s
```
