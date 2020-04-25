include .env

SHORTSHA=`git rev-parse --short HEAD`

LDFLAGS=-X main.appVersion=${VERSION}
LDFLAGS+=-X main.shortSha=$(SHORTSHA)

info:
	@echo Version = ${VERSION}
	@echo Username = ${USERNAME}
	@echo Repository = ${REPONAME}
	@echo Github Token = ${GITHUB_TOKEN}
	@echo Docker Token = ${DOCKER_TOKEN}
	@echo LDFLAGS = $(LDFLAGS)

build:
	@echo Build application ...
	@go build -ldflags "$(LDFLAGS)" .

utils:
	@echo Install tools ...
	@go get github.com/mitchellh/gox
	@go get github.com/tcnksm/ghr

deploy: utils
	@echo Build packages ...
	@CGO_ENABLED=0 gox -os="linux" -arch="amd64" -parallel=4 -ldflags "$(LDFLAGS)" -output "dist/arvand-exporter_{{.OS}}_{{.Arch}}"
	@echo Create release ...
	@ghr -t ${GITHUB_TOKEN} -u ${USERNAME} -r ${REPONAME} -replace ${VERSION} dist/
	@echo Done !

dockerhub:
	@docker login -u ${USERNAME} -p ${DOCKER_TOKEN}
	docker build -t ${USERNAME}/${REPONAME}:${VERSION} .
	docker push ${USERNAME}/${REPONAME}:${VERSION}
