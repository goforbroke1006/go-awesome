all: prepare test
.PHONY: all

prepare:
	go mod download
	go mod tidy
.PHONY: prepare

test:
	go test -short ./...
.PHONY: test
