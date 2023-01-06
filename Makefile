.PHONY: build
build:
	go build -o bot cmd/main.go

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