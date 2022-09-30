.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build



swagger:
	swag init -g ./cmd/app/main.go -o ./doc