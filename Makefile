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
ifeq ($(OS),Windows_NT)
	go env -w CGO_ENABLED=1
	go build -o bin/ ./...
else
	CGO_ENABLED=1 go build -o bin/ ./...
endif

## Remove previous build
clean:
	go clean
ifeq ($(OS),Windows_NT)
	rmdir /s /q bin
else
	rm -rf bin/
endif