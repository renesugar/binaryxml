FROM golang:1.8

WORKDIR /go/src

# Download Go XUnit support
RUN go get -u github.com/jstemmer/go-junit-report

# Add library sources
COPY . github.com/BixData/binaryxml

# Download Go dependencies
WORKDIR github.com/BixData/binaryxml
RUN make gogets

# Run tests
RUN go test -v ./... | $GOPATH/bin/go-junit-report > /test-report.xml
