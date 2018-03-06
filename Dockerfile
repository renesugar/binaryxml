FROM golang:1.8

WORKDIR /go/src

# Download Go XUnit support
RUN go get -u github.com/jstemmer/go-junit-report

# Add library sources
COPY . github.com/BixData/binaryxml

# Download Go dependencies
WORKDIR github.com/BixData/binaryxml
RUN make dependencies

# Build
RUN go install

# Configure runtime environment
VOLUME /go/src/github.com/BixData/binaryxml/target
