.DEFAULT_GOAL := offline

.PHONY: fmt vet build clean offline

## Run go fmt against code
fmt:
	go fmt ./...

## Run go vet against code
vet: fmt
	go vet ./...

## Build the Linux-compatible binary for local serverless testing
build: vet
	set GOOS=linux& set GOARCH=amd64& set CGO_ENABLED=0& go build -o bin/go-rest-test ./cmd/server

## Clean up previous builds
clean:
	rm -rf bin/

## Start local DynamoDB
dynamo:
	docker run -d --name server-dynamodb-container -p 8000:8000 amazon/dynamodb-local

## Run the application locally using Serverless Offline
offline: build
	set ENV=local& serverless offline start --stage local --debug
