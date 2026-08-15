# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
BINARY_NAME=morego
BINARY_DIR=bin

.PHONY: all build test lint fmt clean proto init

all: init build test

init:
	@echo "Initializing workspace dependencies..."
	$(GOCMD) work sync

build:
	@echo "Building CLI binary..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/main.go

test:
	@echo "Running unit tests..."
	$(GOTEST) -v ./cmd/... ./core/... ./drivers/... ./pkg/...

fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

lint:
	@echo "Vetting code..."
	$(GOVET) ./cmd/... ./core/... ./drivers/... ./pkg/...

clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

PROTOC_BIN=protoc
ifeq ($(OS),Windows_NT)
    PROTOC_EXE=tools/protoc/bin/protoc.exe
    ifneq ($(wildcard $(PROTOC_EXE)),)
        PROTOC_BIN=$(PROTOC_EXE)
    endif
else
    PROTOC_SH=tools/protoc/bin/protoc
    ifneq ($(wildcard $(PROTOC_SH)),)
        PROTOC_BIN=$(PROTOC_SH)
    endif
endif

proto:
	@echo "Generating Protocol Buffer bindings (gRPC) using $(PROTOC_BIN)..."
	$(PROTOC_BIN) --proto_path=. --go_out=proto --go_opt=module=github.com/ranjith035/morego/proto --go-grpc_out=proto --go-grpc_opt=module=github.com/ranjith035/morego/proto proto/*.proto
