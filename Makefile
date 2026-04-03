.PHONY: help build install test clean run example

help:
	@echo "Available commands:"
	@echo "  make build      - Build the confg binary"
	@echo "  make install    - Install confg to GOPATH/bin"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make run        - Run example generation"
	@echo "  make example    - Generate config from example YAML"

build:
	@echo "Building confg..."
	go build -o bin/confg ./cmd/confg

install:
	@echo "Installing confg..."
	go install ./cmd/confg

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f examples/basic/config_gen.go

run: build
	@echo "Running confg..."
	./bin/confg --help

example: build
	@echo "Generating config from examples/basic/config.yaml..."
	./bin/confg --path=examples/basic/config.yaml --out=examples/basic/config_gen.go --package=main

.DEFAULT_GOAL := help