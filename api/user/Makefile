GOARCH              ?= amd64
GOOS                ?= linux
VERSION             ?= SNAPSHOT
ENV                 ?= local
ASSETS              := config
SERVICE_NAME        := user
BINARY_NAME         := $(SERVICE_NAME)-$(GOOS)-$(GOARCH)-$(VERSION)
TARBALL_NAME        := $(BINARY_NAME).tar.gz
ARTIFACTS_BUCKET    := childrenofukiyo-artifacts
BUILD_DIR           := build
OUTPUT 				:= main

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f $(OUTPUT) $(PACKAGED_TEMPLATE)

.PHONY: install
install:
	go get ./...

main: main.go
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)/bin
	go build -o $(BUILD_DIR)/bin/$(OUTPUT) main.go

# compile the code to run in Lambda (local or real)
.PHONY: lambda
lambda:
	GOOS=linux GOARCH=amd64 $(MAKE) main

.PHONY: build
build: clean lambda

.PHONY: api
api: build
	doppler run -- sam local start-api -p 8080