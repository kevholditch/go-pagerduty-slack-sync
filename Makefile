
TAG ?= $$(git describe --tags)

build: pagerduty-slack-sync

pagerduty-slack-sync:
	@env GOMODULE111=on CGO_ENABLED=0 go build ./cmd/pagerduty-slack-sync/

vet:
	@go vet -v ./...

check: vet

test:
	@go test -v ./...

docker-build:
	@docker build -t kevholditch/pagerduty-slack-sync:${TAG} -f build/package/Dockerfile .

docker-publish:
	@docker login

ci: build check test

.PHONY: build vet check test