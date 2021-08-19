.PHONY: clean test-coverage test

build: bin go.mod go.sum *.go cmd/midimonster-controller/*.go
	go get
	go build -o bin/midimonster-controller cmd/midimonster-controller/*.go

bin:
	mkdir -p bin

clean:
	rm -rf bin

test:
	go test -cover

test-coverage: bin
	go test -coverprofile=bin/coverage.out
	go tool cover -html bin/coverage.out