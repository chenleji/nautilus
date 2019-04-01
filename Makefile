PROJECT       := nautilus
VERSION       := v0.1
REGISTRY      := hub.docker.com
REGISTRY_USER := $(shell echo ${REGISTRY_USER})
REGISTRY_PWD  := $(shell echo ${REGISTRY_PWD})
REPOSITORY    := ycloud
CI_PIPELINE_ID:= $(shell echo ${CI_PIPELINE_ID})
SONAR_ADDR    := http://10.0.91.242:9000
PKG_LIST      := $(shell go list ./... | grep -v /vendor/)
GO_FILES      := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

test:
	@go test -v ${PKG_LIST}

race: ## Run data race detector
	@go test -race -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	./scripts/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	./scripts/coverage.sh html;

sonar:
	./scripts/sonar.sh ${PROJECT} ${VERSION} ${SONAR_ADDR};

default:
	env GOARCH=amd64 GOOS=linux go build

build:
	docker build -t ${REGISTRY}/${REPOSITORY}/${PROJECT}:${VERSION} -f Dockerfile .

push:
	docker login -u ${REGISTRY_USER} -p ${REGISTRY_PWD} http://${REGISTRY}
	docker push ${REGISTRY}/${REPOSITORY}/${PROJECT}:${VERSION}

release: default build push

build_debug:
	docker build -t ${REGISTRY}/${REPOSITORY}/${PROJECT}:${CI_PIPELINE_ID} -f Dockerfile .

push_debug:
	docker login -u ${REGISTRY_USER} -p ${REGISTRY_PWD} http://${REGISTRY}
	docker push ${REGISTRY}/${REPOSITORY}/${PROJECT}:${CI_PIPELINE_ID}

all: default build_debug push_debug

docker_build:
	docker build -t ${REGISTRY}/${REPOSITORY}/${PROJECT}:latest -f Dockerfile.build .

deploy:
	./scripts/deploy.sh ${PROJECT} ${REGISTRY}/${REPOSITORY}/${PROJECT}:${CI_PIPELINE_ID}
