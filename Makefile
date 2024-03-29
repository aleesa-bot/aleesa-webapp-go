#!/usr/bin/env gmake -f

BUILDOPTS=-ldflags='-s -w' -a -gcflags=all=-l -trimpath

all: clean build

build:
	CGO_ENABLED=0 go build ${BUILDOPTS} -o aleesa-webapp-go \
		types.go globals.go pcache-db-util.go aleesa-webapp-lib.go main.go \
		xkcdru.go randomFox.go theCatAPI.go bunicomic.go \
		anekdotru.go monkeyUser.go obutts.go oboobs.go \
		openweathermap.go prazdnikisegodnyaru.go

clean:
	go clean

upgrade:
	$(RM) -rf vendor
	go get -d -u -t ./...
	go mod tidy
	go mod vendor

# vim: set ft=make noet ai ts=4 sw=4 sts=4:
