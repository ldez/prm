.PHONY: all

GOFILES := $(shell go list -f '{{range $$index, $$element := .GoFiles}}{{$$.Dir}}/{{$$element}}{{"\n"}}{{end}}' ./... | grep -v '/vendor/')
TXT_FILES := $(shell find * -type f -not -path 'vendor/**')

default: clean misspell checks test build-crossbinary

test: clean
	go test -v -cover ./...

dependencies:
	dep ensure -v

clean:
	rm -rf dist/ cover.out

checks: check-fmt
	gometalinter ./...

check-fmt: SHELL := /bin/bash
check-fmt:
	diff -u <(echo -n) <(gofmt -d $(GOFILES))

misspell:
	misspell -source=text -error $(TXT_FILES)

build: clean misspell checks test
	go build

build-crossbinary:
	./.script/crossbinary