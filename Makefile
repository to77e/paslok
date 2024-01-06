include .env.example
export

LOCAL_BIN:=$(CURDIR)/bin
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint

.PHONY: run
run:
	go run cmd/paslok/main.go

.PHONY: build
build:
	go build -ldflags "-X main.version=`git tag --sort=-version:refname | head -n 1`" \
		-o $(GOPATH)/bin/paslok cmd/paslok/main.go

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: install-lint
install-lint:
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Downloading golangci-lint latest)
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
endif

.PHONY: lint
lint: install-lint
	$(GOLANGCI_BIN) run --config=.golangci.yaml ./... -v



.PHONY: local-up
local-up:
	docker-compose -f docker-compose.yaml up -d

.PHONY: local-down
local-down:
	docker-compose -f docker-compose.yaml down -v