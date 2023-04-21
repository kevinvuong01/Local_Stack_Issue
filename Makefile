JANUS_STACK_NAME=whatever
LDFLAGS='-X github.com/cambiahealth/${JANUS_STACK_NAME}/config.buildVersion='"${JANUS_STACK_VERSION}"' -X github.com/cambiahealth/${JANUS_STACK_NAME}/config.buildServiceName='"${JANUS_STACK_NAME}"''
JANUS_STACK_VERSION=$$(if [ -d .git ]; then git rev-parse --short HEAD; else echo "v0.0.0"; fi)
BUILD_ARGS='GOARCH=amd64 CGO_ENABLED=0 GOOS=linux'

up: ## Make and start the app in containers
	@JANUS_STACK_NAME=${JANUS_STACK_NAME} \
  	JANUS_STACK_VERSION=${JANUS_STACK_VERSION} \
  	GITHUB_TOKEN=${GITHUB_TOKEN} \
	docker-compose -f docker-compose.yml up --build

up-local-dev: ## start the local containers except the endpoint one
	docker-compose -f docker-compose.yml up --build

down: ## Stop the running containers
	docker-compose -f docker-compose.yml down


build: ## Build go binary
	@go build \
	-ldflags='-X github.com/cambiahealth/prime-integration-service/config.buildVersion='"${JANUS_STACK_VERSION}"' -X github.com/cambiahealth/prime-integration-service/config.buildServiceName='"${JANUS_STACK_NAME}"'' \
	-o bin/service .


.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: api build clean cuke docker run test test-short test-integration run-command-with-docker-compose-dev go-install-goswagger-ui
