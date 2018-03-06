.PHONY: check
check:
	-mkdir -p target
	go test -v ./... | tee /dev/stderr | $$GOPATH/bin/go-junit-report > target/test-report.xml


.PHONY: dependencies
dependencies:
	go get -u github.com/jstemmer/go-junit-report
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: clean
clean:
	-rm -rf target
