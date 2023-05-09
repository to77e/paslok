LOCAL_BIN:=$(CURDIR)/bin
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
GOLANGCI_TAG:=1.51.2

.PHONY: build
build:
	go build -o bin/pswrd cmd/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: create
create:
	go run cmd/main.go -c $(name) $(comment) -r $(name)

.PHONY: read
read:
	go run cmd/main.go -r $(name)

.PHONY: list
list:
	go run cmd/main.go -l


# install golangci-lint binary
.PHONY: install-lint
install-lint:
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Downloading golangci-lint v$(GOLANGCI_TAG))
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_TAG)
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
endif

# run diff lint like in pipeline
.PHONY: lint
lint: install-lint
	$(info Running lint...)
	$(GOLANGCI_BIN) run --new-from-rev=origin/master --config=.golangci.yaml ./...


# run full lint like in pipeline
.PHONY: lint-full
lint-full: install-lint
	$(GOLANGCI_BIN) run --config=.golangci.yaml ./... -v
