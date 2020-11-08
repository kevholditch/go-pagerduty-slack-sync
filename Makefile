


build:
	@env GOMODULE111=on find ./cmd/* -maxdepth 1 -type d -exec go build "{}" \;

install-lint:
	@go get -u golang.org/x/lint/golint

install-deps: install-lint

lint:
	@golint ./...

vet:
	@go vet -v ./...

check: lint vet

test:
	@go test -v ./...

ci: install-deps build check test

.PHONY: build install-deps install-lint lint vet check test