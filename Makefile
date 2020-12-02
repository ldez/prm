.PHONY: clean check test build

export GO111MODULE=on

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%R:%M:%S%p')

default: clean check test build

clean:
	rm -rf dist/ cover.out

test: clean
	go test -v -cover ./...

build: clean
	@echo Version: $(VERSION) $(BUILD_DATE)
	go build -v -ldflags '-X "github.com/ldez/prm/v3/meta.version=${VERSION}" -X "github.com/ldez/prm/v3/meta.commit=${SHA}" -X "github.com/ldez/prm/v3/meta.date=${BUILD_DATE}"' -trimpath

check:
	golangci-lint run

## Docs
.PHONY: docs-build docs-serve docs-themes

docs-serve:
	@make -C ./docs hugo

docs-build:
	@make -C ./docs hugo-build

docs-themes:
	@make -C ./docs hugo-themes
