# SPDX-License-Identifier: MIT
#
# Copyright (C) 2018-2019 Nadim Kobeissi <nadim@nadim.computer>. All Rights Reserved.

PREFIX ?= /usr/local

all: deps
	@go build -ldflags="-s -w" -o build/diskgem github.com/nadimkobeissiobeissi/diskgem/...

freebsdamd64: deps
	@GOOS="freebsd" GOARCH="amd64" go build -ldflags="-s -w" -o build/diskgem_freebsd_amd64 github.com/nadimkobeissi/diskgem/...

linuxamd64: deps
	@GOOS="linux" GOARCH="amd64" go build -ldflags="-s -w" -o build/diskgem_linux_amd64 github.com/nadimkobeissi/diskgem/...

linuxarm: deps
	@GOOS="linux" GOARCH="arm" go build -ldflags="-s -w" -o build/diskgem_linux_arm github.com/nadimkobeissi/diskgem/...

darwinamd64: deps
	@GOOS="darwin" GOARCH="amd64" go build -ldflags="-s -w" -o build/diskgem_darwin_amd64 github.com/nadimkobeissi/diskgem/...

deps:
	@go get -u ./...

install:
	@install -m0755 build/diskgem $(PREFIX)/bin/diskgem
	@install -m0644 docs/man/diskgem.1 $(PREFIX)/share/man/man1/diskgem.1
ifeq ($(shell uname),Darwin)
	@/usr/libexec/makewhatis	
else
	@mandb -q
endif

lint:
	golangci-lint run

clean:
	@$(RM) -rf build/diskgem*

.PHONY: all freebsdamd64 linuxamd64 linuxarm darwinamd64 deps install lint clean docs build cmd
