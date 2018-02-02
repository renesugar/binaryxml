#!/bin/bash
go test -v ./... | tee /dev/tty | $GOPATH/bin/go-junit-report > target/test-report.xml
