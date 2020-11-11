
TAG ?= $$(git describe --tags)

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

docker-build:
	@echo ${TAG}
	@docker build -t kevholditch/pagerduty-slack-sync:${TAG} -f build/package/Dockerfile .

ci: install-deps build check test docker-build

.PHONY: build install-deps install-lint lint vet check test