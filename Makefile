.PHONY: check
check:
	-mkdir -p target
	go test -v ./... | tee /dev/stderr | $$GOPATH/bin/go-junit-report > target/test-report.xml

  
.PHONY: gogets
gogets:
	go get github.com/cevaris/ordered_map
	go get github.com/docktermj/go-logger/logger
	go get github.com/jnewmoyer/xmlpath
	go get github.com/stretchr/testify/assert
	go get github.com/tdewolff/minify/xml


.PHONY: clean
clean:
	-rm -rf target
