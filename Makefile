.PHONY: all

PKGS := $(shell go list ./... | grep -v '/vendor/')
GOFILES := $(shell go list -f '{{range $$index, $$element := .GoFiles}}{{$$.Dir}}/{{$$element}}{{"\n"}}{{end}}' ./... | grep -v '/vendor/')
TXT_FILES := $(shell find * -type f -not -path 'vendor/**')

default: clean misspell checks lint test build-crossbinary

test: clean
	go test -v -cover ./...

dependencies:
	dep ensure -v

clean:
	rm -f cover.out

lint:
	golint -set_exit_status $(PKGS)

checks: check-fmt
	staticcheck $(PKGS)
	gosimple $(PKGS)

check-fmt: SHELL := /bin/bash
check-fmt:
	diff -u <(echo -n) <(gofmt -d $(GOFILES))

misspell:
	misspell -source=text -error $(TXT_FILES)

build: clean misspell checks lint test
	go build

build-crossbinary:
	./.script/crossbinary