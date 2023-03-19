.PHONY: clean test-coverage test
PACKAGE_NAME := github.com/jduepmeier/midimonster-controller

build: bin go.mod go.sum *.go cmd/midimonster-controller/*.go
	go get
	go build -o bin/midimonster-controller cmd/midimonster-controller/*.go


build-arm64:
	GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 go build -o bin/midimonster-controller-arm64 cmd/midimonster-controller/*.go
build-armv7:
	GOOS=linux GOARCH=arm GOARM=7 CC=arm-none-eabi-gcc CGO_ENABLED=1 go build -o bin/midimonster-controller-arm64 cmd/midimonster-controller/*.go


bin:
	mkdir -p bin

clean:
	rm -rf bin

test:
	go test -cover

release-container:
	podman run \
		--rm \
		-e GITHUB_TOKEN \
		-v /usr/include/systemd:/usr/include/systemd \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:v1.20.2 \
		release --clean

test-coverage: bin
	go test -coverprofile=bin/coverage.out
	go tool cover -html bin/coverage.out
