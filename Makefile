# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMODTIDY=$(GOCMD) mod tidy
GOMODVERIFY=$(GOCMD) mod verify
BINARY_NAME=goproject
BINARY_OSX=$(BINARY_NAME)_darwin
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_FREEBSD=$(BINARY_NAME)_freebsd
DIR_NAME := $(shell pwd)

# Full build commands
all: clean mods format lint test build
release: mods format lint clean test build-osx build-linux build-freebsd package

# Code cleanup commands
mods:
	$(GOMODTIDY)
	$(GOMODVERIFY)
format:
	find . -name \*.go -not -path vendor -not -path target -exec $(GOPATH)/bin/goimports -w {} \;
lint:
	${GOPATH}/bin/golangci-lint run --enable-all --disable funlen --disable gochecknoglobals --disable gochecknoinits

# Build commands
build: 
	$(GOBUILD) -gcflags="-trimpath=$(DIR_NAME)" -o target/$(BINARY_NAME) -v cmd/$(BINARY_NAME)/main.go
build-debug: 
	$(GOBUILD) -gcflags="-m -l" -ldflags="-v" -o target/$(BINARY_NAME) -v cmd/$(BINARY_NAME)/main.go
test: 
	$(GOTEST) -v ./...
clean: 
	rm -f target/*
run:
	$(GOBUILD) -gcflags="-trimpath=$(DIR_NAME)" -o target/$(BINARY_NAME) -v cmd/$(BINARY_NAME)/main.go
	./target/$(BINARY_NAME)

# Cross compilation
build-osx:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -gcflags="-trimpath=$(DIR_NAME)" -o target/$(BINARY_OSX) -v cmd/$(BINARY_NAME)/main.go
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -gcflags="-trimpath=$(DIR_NAME)" -o target/$(BINARY_LINUX) -v cmd/$(BINARY_NAME)/main.go
build-freebsd:
	CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 $(GOBUILD) -gcflags="-trimpath=$(DIR_NAME)" -o target/$(BINARY_FREEBSD) -v cmd/$(BINARY_NAME)/main.go

# Packaging
package: package-osx package-linux package-freebsd
package-osx:
	tar -zcvf target/$(BINARY_OSX).tar.gz target/$(BINARY_OSX)
package-linux:
	tar -zcvf target/$(BINARY_LINUX).tar.gz target/$(BINARY_LINUX)
package-freebsd:
	tar -zcvf target/$(BINARY_FREEBSD).tar.gz target/$(BINARY_FREEBSD)