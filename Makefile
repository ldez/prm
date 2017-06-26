.PHONY: all

default:
	./.script/make.sh

dependencies:
	glide install

build:
	go build

validate:
	./.script/make.sh validate-glide validate-gofmt validate-govet validate-golint validate-misspell

test-unit:
	./.script/make.sh test-unit
