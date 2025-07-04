all: test

.PHONY: all generate-openapi test lint

generate-openapi: generate-swagger rename-swagger

generate-swagger:
	swag init --outputTypes json

rename-swagger:
	mv ./docs/swagger.json ./docs/users-openapi.json

test:
	go test ./... -coverprofile=coverage.out -coverpkg=./...

lint: install-lint run-lint

install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.2.1

run-lint:
	golangci-lint run
