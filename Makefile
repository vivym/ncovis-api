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

REGISTRY_URL = registry.cn-beijing.aliyuncs.com/viv/
IMAGE_NAME = ncovis-api
IMAGE_VER = ${VERSION}-${COMMIT_HASH}
IMAGE_FULL_NAME = ${REGISTRY_URL}${IMAGE_NAME}:${IMAGE_VER}

.PHONY: all
all: build

.PHONY: clean
clean:
	rm -f build/ncovis-server

build: clean
	go build ${BUILD_PARAMS}

.PHONY: build
release: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${BUILD_PARAMS}
	docker build -t ${IMAGE_NAME}:${IMAGE_VER} ./build
	docker tag ${IMAGE_NAME}:${IMAGE_VER} ${IMAGE_FULL_NAME}
	docker push ${IMAGE_FULL_NAME}
	sed 's/__IMAGE_FULL_NAME__/${IMAGE_FULL_NAME}/g' deployment.yaml > build/deployment.yaml
