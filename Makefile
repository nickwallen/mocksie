all: test

lint:
	golangci-lint run ./...

test: lint
	go test ./...

install:
	go install ./cmd/mocksie