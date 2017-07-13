# Binary XML Library

This GoLang library provides encoding and decoding support for NuBix binary xml.

## Usage

### Convert Binary XML to XML

```go
import (
	"github.com/BixData/binaryxml"
	"io/ioutil"
)

binaryXml, _ := ioutil.ReadFile("mydata.binaryxml")
xml, err := binaryxml.ToXML(binaryXml)
```

### Encode a struct

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

### Routing requests

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
		MOID          string          `xml:"moid"`
		MID           string          `xml:"mid"`
		Data          bixResponseData `xml:"Data"`
	}
	type bixResponseData struct {
		XMLName struct{} `xml:"BixResponse"`
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

```
$ mkdir -p workspaces/nubix/agent/src/github.com/BixData
$ cd workspaces/nubix/agent
$ export GOPATH=`pwd`
$ cd src/github.com/BixData
$ git clone <this repo>
$ cd <this repo>
$ make gogets
```

And then test with:

```
$ go test
...
PASS
ok  	github.com/BixData/binaryxml	0.036s
```
