.PHONY: clean check test build dependencies fmt imports hugo-theme hugo-theme-clean hugo-build hugo

export GO111MODULE=on

GOFILES := $(shell git ls-files '*.go' | grep -v '^vendor/')

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

default: clean check test build

dependencies:
	go mod download

clean:
	rm -rf dist/ cover.out

test: clean
	go test -v -cover ./...

build: clean
	@echo Version: $(VERSION) $(BUILD_DATE)
	go build -v -ldflags '-X "github.com/ldez/prm/meta.version=${VERSION}" -X "github.com/ldez/prm/meta.commit=${SHA}" -X "github.com/ldez/prm/meta.date=${BUILD_DATE}"'

check:
	golangci-lint run

fmt:
	@gofmt -s -l -w $(GOFILES)

imports:
	@goimports -w $(GOFILES)

## Docs
.PHONY: docs-build docs-serve docs-themes

docs-serve:
	@make -C ./docs hugo

docs-build:
	@make -C ./docs hugo-build

docs-themes:
	@make -C ./docs hugo-themes
