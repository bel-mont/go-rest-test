.DEFAULT_GOAL := build

.PHONY: fmt vet build bin
## Run go fmt against code
fmt:
	go fmt ./...

## Run go vet against code
vet: fmt
	go vet ./...

## Build the binary file
## Remove if sqlite3 is removed. And then set CGO_ENABLED=0
build: vet
	go build -o bin/go-rest-test.exe ./cmd/server

## Remove previous build
clean:
	go clean
ifeq ($(OS),Windows_NT)
	rmdir /s /q bin
else
	rm -rf bin/
endif