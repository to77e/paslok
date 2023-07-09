include .env
export

LOCAL_BIN:=$(CURDIR)/bin
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint

.PHONY: build
build:
	 go build -o $(GOPATH)/bin/paslok cmd/paslok/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: create
create:
	go run cmd/paslok/main.go -c $(name) $(comment) -r $(name)

.PHONY: read
read:
	go run cmd/paslok/main.go -r $(name)

.PHONY: list
list:
	go run cmd/paslok/main.go -l

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
