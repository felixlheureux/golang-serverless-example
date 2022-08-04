GOARCH              ?= arm64
GOOS                ?= linux
VERSION             ?= SNAPSHOT
ENV                 ?= local
ASSETS              := config
SERVICE_NAME        := odin
BINARY_NAME         := $(SERVICE_NAME)-$(GOOS)-$(GOARCH)-$(VERSION)
TARBALL_NAME        := $(BINARY_NAME).tar.gz
ARTIFACTS_BUCKET    := childrenofukiyo-artifacts
BUILD_DIR           := build

test:
	(cd ./user && make test)
	(cd ./auth && make test)

build: # build a distribution tarball
	# clean build directory
	rm -rf $(BUILD_DIR)
	# create build directory
	mkdir -p $(BUILD_DIR)/bin
	# build all services
	(cd ./user && make build)
	(cd ./auth && make build)

publish:
	(cd ./user && make publish)
	(cd ./auth && make publish)

all:
	(cd ./user && make all)
	(cd ./auth && make all)

clean:
	rm -rf $(BUILD_DIR)

.PHONY: api
api:
	doppler run -- sam local start-api -p 8080 --skip-pull-image