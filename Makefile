# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=goproject
BINARY_LINUX=$(BINARY_NAME)_linux

all: test build
build: 
	$(GOBUILD) -o target/$(BINARY_NAME) -v cmd/$(BINARY_NAME)/main.go
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_LINUX)
run:
	$(GOBUILD) -o target/$(BINARY_NAME) -v cmd/$(BINARY_NAME)/main.go
	./target/$(BINARY_NAME)

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o target/$(BINARY_NAME) -v cmd/$(BINARY_NAME)/main.go
