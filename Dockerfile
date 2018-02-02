FROM golang:1.8

WORKDIR /go/src

# Download Go XUnit support
RUN go get -u github.com/jstemmer/go-junit-report

# Add library sources
COPY . github.com/BixData/binaryxml

# Download Go dependencies
WORKDIR github.com/BixData/binaryxml
RUN make gogets

# Build
RUN go install

# Configure runtime environment
ADD runtests.sh runtests.sh
VOLUME /go/src/github.com/BixData/binaryxml/target
