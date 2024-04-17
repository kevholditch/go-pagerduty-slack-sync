DOCKER_COMPOSE_FILE=deployments/local/docker-compose.yaml
TAG ?= $$(git describe --tags)

build:
	@env GOMODULE111=on find ./cmd/* -maxdepth 1 -type d -exec go build "{}" \;

vet:
	@go vet -v ./...

check: vet

test:
	@go test -v ./...

docker-build:
	@docker build -t pagerduty-slack-sync:${TAG} -f build/package/Dockerfile .

docker-publish:
	@docker login

dev:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up --build pagerduty-slack-sync

ci: build check test

.PHONY: build vet check test