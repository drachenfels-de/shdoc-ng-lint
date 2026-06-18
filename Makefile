GOBIN ?= $(CURDIR)/.tools/bin
GOCACHE ?= $(CURDIR)/.tools/gocache
GOMODCACHE ?= $(CURDIR)/.tools/gomodcache
PATH := $(GOBIN):$(CURDIR)/.tools/go/bin:$(PATH)

.PHONY: tools fmt lint test build ci

tools:
	GOBIN=$(GOBIN) GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GO) install mvdan.cc/gofumpt@v0.8.0
	GOBIN=$(GOBIN) GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2

fmt:
	find . -path './.tools' -prune -o -name '*.go' -print | xargs $(GOBIN)/gofumpt -w

lint:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GOBIN)/golangci-lint run ./...

test:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go test ./...

build:
	CGO_ENABLED=0 GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go build -ldflags="-w -s" .

ci: fmt lint test build

release:
	./release.sh

