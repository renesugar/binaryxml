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
	XMLName struct{} `xml:"BixRequest"`
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

## Testing

Setup a workspace:

```
$ mkdir -p workspaces/nubix/agent/src/github.com/BixData
$ cd workspaces/nubix
$ export GOPATH=`pwd`
$ cd agent/src/github.com/BixData
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
