.PHONY: build ui

GO111MODULE ?= on

.EXPORT_ALL_VARIABLES:

all: fmt deps build

deps:
	go generate -x
	go get ./...

fmt:
	gofmt -w .
	go vet .

build: fmt
	go build -o bin/autotag .
	which autotag && cp bin/autotag $(shell which autotag)
