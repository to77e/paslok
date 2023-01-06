.PHONY: run 
run:
	go run cmd/main.go

.PHONY: build
build:
	go build -o bot cmd/main.go

.PHONY: test
test:
	go test -v ./...