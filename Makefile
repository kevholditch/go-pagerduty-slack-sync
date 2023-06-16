
TAG ?= $$(git describe --tags)

build:
	go build

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

docker-build:
	@docker build -t kevholditch/pagerduty-slack-sync:${TAG} -f build/package/Dockerfile .

docker-publish:
	@docker login

ci: install-deps build check test

.PHONY: build install-deps install-lint lint vet check test