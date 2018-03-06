FROM golang:1.9

ARG REFRESHED_AT=2018-03-05

# ============================================================
# Add sources and install dependencies
# ============================================================

COPY . $GOPATH/src/github.com/BixData/binaryxml
WORKDIR $GOPATH/src/github.com/BixData/binaryxml
RUN make dependencies

# ============================================================
# Build
# ============================================================

WORKDIR $GOPATH/src/github.com/BixData/binaryxml
ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go install $GO_PACKAGE

# Configure runtime environment
VOLUME /go/src/github.com/BixData/binaryxml/target
