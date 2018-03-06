FROM golang:1.9

ARG REFRESHED_AT=2018-03-05

WORKDIR /go/src

# Add library sources
COPY . github.com/BixData/binaryxml

# Download Go dependencies
WORKDIR github.com/BixData/binaryxml
RUN make dependencies

# Build
ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go install

# Configure runtime environment
VOLUME /go/src/github.com/BixData/binaryxml/target
