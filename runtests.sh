#!/bin/bash
go test -v ./... | tee /dev/stderr | $GOPATH/bin/go-junit-report > target/test-report.xml
