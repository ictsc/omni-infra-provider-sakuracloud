.PHONY: build test fmt check

build:
	mkdir -p bin
	go build -o bin/omni-infra-provider-sakuracloud ./cmd/omni-infra-provider-sakuracloud

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -type f)

check:
	test -z "$$(gofmt -l $$(find . -name '*.go' -type f))"
	go vet ./...
	go test ./...
