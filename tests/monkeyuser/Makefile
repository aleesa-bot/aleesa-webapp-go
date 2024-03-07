#!/usr/bin/env gmake -f

BUILDOPTS=-ldflags='-s -w' -a -gcflags=all=-l -trimpath
BINARY=monkeyuser-test

all: clean build

build:
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${BINARY} aleesa-webapp-lib.go  globals.go  main.go  monkeyUser.go  types.go

clean:
	$(RM) ${BINARY}

upgrade:
	go get -d -u -t ./...
	go mod tidy

# vim: set ft=make noet ai ts=4 sw=4 sts=4:
