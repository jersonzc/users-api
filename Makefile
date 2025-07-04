all: test

.PHONY: all test

test:
	go test ./... -coverprofile=coverage.out -coverpkg=./...

lint: install-lint run-lint

install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.2.1

run-lint:
	golangci-lint run
