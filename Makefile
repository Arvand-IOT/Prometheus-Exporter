include .env
include .version

SHORTSHA=`git rev-parse --short HEAD`
LDFLAGS=-X main.appVersion=${VERSION}
LDFLAGS+=-X main.shortSha=$(SHORTSHA)
GOCMD=go
GOTEST=$(GOCMD) test

.PHONY: info dev build utils docker docker-dev help
.DEFAULT_GOAL := help


info: ## Show all variables
	@echo Version = ${VERSION}
	@echo Repository = ${REPONAME}
	@echo Username = ${DOCKER_USERNAME}
	@echo Docker Token = ${DOCKER_TOKEN}
	@echo LDFLAGS = $(LDFLAGS)

dev: ## Build the staging binaries
	$(GOCMD) build -ldflags "$(LDFLAGS)" -o "dist/arvand-exporter"

utils: ## Install required modules
	@$(GOCMD) install github.com/mitchellh/gox@latest

build: utils ## Build binary
	CGO_ENABLED=0 gox -os="linux" -arch="amd64 arm64" -parallel=4 -ldflags "$(LDFLAGS)" -output "dist/arvand-exporter_{{.OS}}_{{.Arch}}"

docker: build ## Build & push Docker image
	@echo ${DOCKER_TOKEN} | docker login --username ${DOCKER_USERNAME} --password-stdin
	docker buildx create --name mybuilder --use
	docker buildx build -t ${DOCKER_USERNAME}/${REPONAME}:${VERSION} -t ${DOCKER_USERNAME}/${REPONAME}:latest --platform=linux/arm64,linux/amd64 . --push

docker-dev: build ## Build Docker image
	docker buildx build -t ${DOCKER_USERNAME}/${REPONAME}:${VERSION} . --load

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
