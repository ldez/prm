.PHONY: all

GOFILES := $(shell go list -f '{{range $$index, $$element := .GoFiles}}{{$$.Dir}}/{{$$element}}{{"\n"}}{{end}}' ./... | grep -v '/vendor/')

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

default: clean checks test build

test: clean
	go test -v -cover ./...

dependencies:
	dep ensure -v

clean:
	rm -rf dist/ cover.out

build: clean
	@echo Version: $(VERSION) $(BUILD_DATE)
	go build -v -ldflags '-X "github.com/ldez/prm/meta.version=${VERSION}" -X "github.com/ldez/prm/meta.commit=${SHA}" -X "github.com/ldez/prm/meta.date=${BUILD_DATE}"'

checks: check-fmt
	golangci-lint run

check-fmt: SHELL := /bin/bash
check-fmt:
	diff -u <(echo -n) <(gofmt -d $(GOFILES))
