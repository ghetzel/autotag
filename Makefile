.PHONY: build ui

all: fmt deps build

deps:
	@go list golang.org/x/tools/cmd/goimports || go get golang.org/x/tools/cmd/goimports
	go generate -x
	go get .

fmt:
	goimports -w .
	go vet .

build: fmt
	go build -o bin/autotag .
	which autotag && cp bin/autotag $(shell which autotag)
