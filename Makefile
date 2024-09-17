#!/usr/bin/env gmake -f

BUILDOPTS=-ldflags='-s -w' -a -gcflags=all=-l -trimpath

all: clean build

build:
	CGO_ENABLED=0 go build ${BUILDOPTS} -o aleesa-webapp-go ./cmd/aleesa-webapp-go

buildutils:
	CGO_ENABLED=0 go build ${BUILDOPTS} -o flickr_init ./cmd/flickr_init
	CGO_ENABLED=0 go build ${BUILDOPTS} -o flickr_populate ./cmd/flickr_populate
	CGO_ENABLED=0 go build ${BUILDOPTS} -o flickr_test ./cmd/flickr_test

clean:
	$(RM) -rf aleesa-webapp-go flickr_init flickr_populate flickr_test

upgrade:
	$(RM) -rf vendor
	go get -d -u -t ./...
	go mod tidy
	go mod vendor

# vim: set ft=make noet ai ts=4 sw=4 sts=4:
