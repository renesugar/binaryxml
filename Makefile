TARGET ?= ./target
TEST_REPORT ?= $(TARGET)/test-report.xml

.PHONY: check
check:
	mkdir -p $(TARGET)
	go test -v ./... | tee /dev/stderr | $$GOPATH/bin/go-junit-report > $(TEST_REPORT)

.PHONY: dependencies
dependencies:
	go get -u github.com/jstemmer/go-junit-report
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: clean
clean:
	-rm -rf $(TARGET)
