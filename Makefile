.PHONY: build test run all install
all: build

install:
	@go get -d ./...

build:
	@go build ./...

test:
	@go test ./...

test-cover:
	@echo "mode: set" > acc.coverage-out
	@go test -coverprofile=main.coverage-out .
	@cat main.coverage-out | grep -v "mode: set" >> acc.coverage-out
	@go tool cover -html=acc.coverage-out
	@rm *.coverage-out

run:
	@go run *.go
