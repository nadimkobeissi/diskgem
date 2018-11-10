# SPDX-License-Identifier: GPL-2.0
#
# Copyright (C) 2018 Nadim Kobeissi <nadim@nadim.computer>. All Rights Reserved.

PREFIX ?= /usr

ifeq ($(shell uname -s),Linux)
OSFLAG=linux
ifeq ($(shell uname -p),x86_64)
ARCHFLAG=amd64
else
ARCHFLAG=arm
endif
endif
ifeq ($(shell uname -s),Darwin)
OSFLAG=darwin
ARCHFLAG=amd64
endif

all: deps
	@export GOOS="$(OSFLAG)"
	@export GOARCH="$(ARCHFLAG)"
	go build -ldflags="-s -w" -o dist/diskgem_$(OSFLAG)_$(ARCHFLAG) src/*

deps:
	@cd src
	@go get -d ./...
	@cd ..

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

install:
	install -m0755 dist/diskgem_$(OSFLAG)_$(ARCHFLAG) $(PREFIX)/bin/diskgem
	install -m644 man/diskgem.1 $(PREFIX)/share/man/man1/diskgem.1
	mandb -q

clean:
	rm -rf dist/diskgem*

.PHONY: all clean man src
