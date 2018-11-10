# SPDX-License-Identifier: MIT
#
# Copyright (C) 2018 Nadim Kobeissi <nadim@nadim.computer>. All Rights Reserved.

PREFIX ?= /usr

all: deps
	go build -ldflags="-s -w" -o dist/diskgem src/*

linuxamd64: deps
	@export GOOS="linux"
	@export GOARCH="amd64"
	go build -ldflags="-s -w" -o dist/diskgem_linux_amd64 src/*

linuxarm: deps
	@export GOOS="linux"
	@export GOARCH="arm"
	go build -ldflags="-s -w" -o dist/diskgem_linux_arm src/*

darwinamd64: deps
	@export GOOS="darwin"
	@export GOARCH="amd64"
	go build -ldflags="-s -w" -o dist/diskgem_darwin_amd64 src/*

deps:
	@cd src
	@go get -d ./...
	@cd ..

install:
	install -m0755 dist/diskgem $(PREFIX)/bin/diskgem
	install -m644 man/diskgem.1 $(PREFIX)/share/man/man1/diskgem.1
	mandb -q

clean:
	rm -rf dist/diskgem*

.PHONY: all clean man src
