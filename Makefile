#!/usr/bin/env gmake -f

GOOPTS=CGO_ENABLED=0
BUILDOPTS=-ldflags="-s -w" -a -gcflags=all=-l -trimpath -buildvcs=false

MYNAME=aleesa-webapp-go
BINARY=${MYNAME}
UNIX_BINARY=${MYNAME}
WINDOWS_BINARY=${MYNAME}.exe

MYNAME2=flickr_init
BINARY2=${MYNAME2}
UNIX_BINARY2=${MYNAME2}
WINDOWS_BINARY2=${MYNAME2}.exe

MYNAME3=flickr_populate
BINARY3=${MYNAME3}
UNIX_BINARY3=${MYNAME3}
WINDOWS_BINARY3=${MYNAME3}.exe

MYNAME4=flickr_test
BINARY4=${MYNAME4}
UNIX_BINARY4=${MYNAME4}
WINDOWS_BINARY4=${MYNAME4}.exe

RMCMD=rm -rf

# На windows имя бинарника может зависеть не только от платформы, но и от выбранной цели, для linux-а суффикс .exe
# не нужен
ifeq ($(OS),Windows_NT)
ifdef GOOS
ifeq ($(GOOS),windows)
BINARY=${WINDOWS_BINARY}
BINARY2=${WINDOWS_BINARY2}
BINARY3=${WINDOWS_BINARY3}
BINARY4=${WINDOWS_BINARY4}
else  # not ifeq ($(GOOS),windows)
BINARY=${MYNAME}
BINARY2=${MYNAME2}
BINARY3=${MYNAME3}
BINARY4=${MYNAME4}
endif # ifeq ($(GOOS),windows)
else  # not ifdef GOOS
BINARY=${WINDOWS_BINARY}
BINARY2=${WINDOWS_BINARY}
BINARY3=${WINDOWS_BINARY}
BINARY4=${WINDOWS_BINARY}
endif # ifdef GOOS
ifeq ($(SHELL), sh.exe)
RMCMD=DEL /Q /F
endif
endif

# Явно определяем символ новой строки, чтобы избежать неоднозначности на windows
define IFS

endef


all: clean build


build:
ifeq ($(OS),Windows_NT)
# Looks like on windows gnu make explicitly set SHELL to sh.exe, if it was not set.
ifeq ($(SHELL), sh.exe)
#       # Vanilla cmd.exe / powershell.
	SET "CGO_ENABLED=0"
	go build ${BUILDOPTS} -o ${BINARY} ./cmd/${MYNAME}
else ifeq (,$(findstring(Git/usr/bin/sh.exe, $(SHELL))))
#       # git-bash
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${BINARY} ./cmd/${MYNAME}
else  # not ifeq (,$(findstring(Git/usr/bin/sh.exe, $(SHELL))))
#       # Some other shell.
#       # TODO: handle it.
	$(info "-- Dunno how to handle this shell: ${SHELL}")
endif # ifeq (,$(findstring(Git/usr/bin/sh.exe, $(SHELL))))
else  # not  ($(OS),Windows_NT)
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${BINARY} ./cmd/${MYNAME}
endif # ifeq ($(OS),Windows_NT)

buildutils:
ifeq ($(OS),Windows_NT)
# Looks like on windows gnu make explicitly set SHELL to sh.exe, if it was not set.
ifeq ($(SHELL), sh.exe)
#       # Vanilla cmd.exe / powershell.
	SET "CGO_ENABLED=0"
	go build ${BUILDOPTS} -o ${BINARY2} ./cmd/${MYNAME2}
	go build ${BUILDOPTS} -o ${BINARY3} ./cmd/${MYNAME3}
	go build ${BUILDOPTS} -o ${BINARY4} ./cmd/${MYNAME4}
else ifeq (,$(findstring(Git/usr/bin/sh.exe, $(SHELL))))
#       # git-bash
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${BINARY2} ./cmd/${MYNAME2}
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${BINARY3} ./cmd/${MYNAME3}
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${BINARY4} ./cmd/${MYNAME4}
else  # not ifeq (,$(findstring(Git/usr/bin/sh.exe, $(SHELL))))
#       # Some other shell.
#       # TODO: handle it.
	$(info "-- Dunno how to handle this shell: ${SHELL}")
endif # ifeq (,$(findstring(Git/usr/bin/sh.exe, $(SHELL))))
else  # not  ($(OS),Windows_NT)
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${BINARY2} ./cmd/${MYNAME2}
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${BINARY3} ./cmd/${MYNAME3}
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${BINARY4} ./cmd/${MYNAME4}
endif # ifeq ($(OS),Windows_NT)


clean:
ifeq ($(OS),Windows_NT)
ifeq ($(SHELL),sh.exe)
#	# Vanilla cmd.exe / powershell.
	if exist ${WINDOWS_BINARY} ${RMCMD} ${WINDOWS_BINARY}
	if exist ${UNIX_BINARY} ${RMCMD} ${UNIX_BINARY}

	if exist ${WINDOWS_BINARY2} ${RMCMD} ${WINDOWS_BINARY2}
	if exist ${UNIX_BINARY2} ${RMCMD} ${UNIX_BINARY2}

	if exist ${WINDOWS_BINARY3} ${RMCMD} ${WINDOWS_BINARY3}
	if exist ${UNIX_BINARY3} ${RMCMD} ${UNIX_BINARY3}

	if exist ${WINDOWS_BINARY4} ${RMCMD} ${WINDOWS_BINARY4}
	if exist ${UNIX_BINARY4} ${RMCMD} ${UNIX_BINARY4}
else  # not ifeq ($(SHELL),sh.exe)
	${RMCMD} ./${WINDOWS_BINARY}
	${RMCMD} ./${UNIX_BINARY}

	${RMCMD} ./${WINDOWS_BINARY2}
	${RMCMD} ./${UNIX_BINARY2}

	${RMCMD} ./${WINDOWS_BINARY3}
	${RMCMD} ./${UNIX_BINARY3}

	${RMCMD} ./${WINDOWS_BINARY4}
	${RMCMD} ./${UNIX_BINARY4}
endif # ifeq ($(SHELL),sh.exe)
else  # not ifeq ($(OS),Windows_NT)
	${RMCMD} ./${BINARY}
	${RMCMD} ./${BINARY2}
	${RMCMD} ./${BINARY3}
	${RMCMD} ./${BINARY4}
endif

upgrade:
	go get -u ./...
	go mod tidy
	go mod vendor

# vim: set ft=make noet ai ts=4 sw=4 sts=4:
