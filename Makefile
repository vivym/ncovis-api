# Build variables
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
DATE_FMT = +%FT%T%z
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif

LDFLAGS += -w -s
LDFLAGS += -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE}

BUILD_PARAMS = -a -ldflags "${LDFLAGS}" -o build/ncovis-server cmd/ncovis-server/*.go

REGISTRY_URL = registry.cn-beijing.aliyuncs.com/public-api/
IMAGE_NAME = ncovis-api
IMAGE_VER = ${VERSION}-${COMMIT_HASH}
IMAGE_FULL_NAME = ${REGISTRY_URL}${IMAGE_NAME}:${IMAGE_VER}

.PHONY: all
all: build

.PHONY: clean
clean:
	@rm -r -f build

.PHONY: build
build: clean
	@mkdir build
	go build ${BUILD_PARAMS}

.PHONY: release
release: clean
	@mkdir build

	@echo
	@echo ---------------------------------------------------------------
	@echo - `date "+%H:%M:%S"` [1] building binary
	@echo ---------------------------------------------------------------
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${BUILD_PARAMS}

	@echo
	@echo ---------------------------------------------------------------
	@echo - `date "+%H:%M:%S"` [2] building docker image
	@echo ---------------------------------------------------------------
	@docker build -t ${IMAGE_NAME}:${IMAGE_VER} -f ./Dockerfile ./build
	@docker tag ${IMAGE_NAME}:${IMAGE_VER} ${IMAGE_FULL_NAME}

	@echo
	@echo ---------------------------------------------------------------
	@echo - `date "+%H:%M:%S"` [3] pushing docker image
	@echo ---------------------------------------------------------------
	@docker push ${IMAGE_FULL_NAME}

	@echo
	@echo ---------------------------------------------------------------
	@echo - `date "+%H:%M:%S"` [4] create k8s deployment.yaml
	@echo ---------------------------------------------------------------
	@sed 's#__IMAGE_FULL_NAME__#${IMAGE_FULL_NAME}#g;s#__GRAPHIQL_TOKEN__#${GRAPHIQL_TOKEN}#g' deployment.yaml > build/deployment.yaml

	@echo
	@echo ---------------------------------------------------------------
	@echo - `date "+%H:%M:%S"` [5] done
	@echo ---------------------------------------------------------------

.PHONY: docker-login
docker-login:
	@echo building $(shell date "$(DATE_FMT)")
	@docker login --username=${REGISTRY_USERNAME} --password=${REGISTRY_PASSWORD} ${REGISTRY_URL}
	@echo done $(shell date "$(DATE_FMT)")
