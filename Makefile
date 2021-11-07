all: lint test

lint:
	golangci-lint run ./...

test:
	go test ./...

install:
	go install ./cmd/mocksie