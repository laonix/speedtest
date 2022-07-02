#!/usr/bin/make

build:
	mkdir -p bin/
	$(shell which go) build -trimpath -o bin/ -v ./cmd/speedtest/...

lint:
	$(shell which golangci-lint) run --verbose --config .golangci.yml

fix-lint:
	$(shell which golangci-lint) run --verbose --config .golangci.yml --fix

test-unit:
	$(shell which go) test -v -race ./...

test-unit-coverage:
	$(shell which go) test -cover ./...

test-benchmark:
	$(shell which go) test -v -bench ./...