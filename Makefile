# SPDX-License-Identifier: MIT
#
# Copyright (C) 2018-2019 Nadim Kobeissi <nadim@nadim.computer>. All Rights Reserved.

PREFIX ?= /usr/local

all: deps
	go build -ldflags="-s -w" -o dist/diskgem src/*

freebsdamd64: deps
	GOOS="freebsd" GOARCH="amd64" go build -ldflags="-s -w" -o dist/diskgem_freebsd_amd64 src/*

linuxamd64: deps
	GOOS="linux" GOARCH="amd64" go build -ldflags="-s -w" -o dist/diskgem_linux_amd64 src/*

linuxarm: deps
	GOOS="linux" GOARCH="arm" go build -ldflags="-s -w" -o dist/diskgem_linux_arm src/*

darwinamd64: deps
	GOOS="darwin" GOARCH="amd64" go build -ldflags="-s -w" -o dist/diskgem_darwin_amd64 src/*

deps:
	@cd src; go get -d ./...; cd ..

install:
	install -m0755 dist/diskgem $(PREFIX)/bin/diskgem
	install -m0644 man/diskgem.1 $(PREFIX)/share/man/man1/diskgem.1
ifeq ($(shell uname),Darwin)
	/usr/libexec/makewhatis	
else
	mandb -q
endif

clean:
	rm -rf dist/diskgem*

.PHONY: all freebsdamd64 linuxamd64 linuxarm darwinamd64 deps install clean src man web
