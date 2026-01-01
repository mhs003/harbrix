PREFIX    ?= /usr/local
BINDIR    ?= $(PREFIX)/bin

GO        ?= go
INSTALL   ?= install

VERSION   := $(shell cat version)

.PHONY: all build install uninstall clean run-daemon run-cli list status start stop

all: build

build:
	mkdir -p build
	$(GO) build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o build/harbrix ./cmd/harbrix
	$(GO) build -trimpath -ldflags="-s -w" -o build/harbrixd ./cmd/harbrixd
	
build-all:
	mkdir -p build/linux
	GOOS=linux GOARCH=amd64   $(GO) build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o build/linux/harbrix-amd64   ./cmd/harbrix
	GOOS=linux GOARCH=amd64   $(GO) build -trimpath -ldflags="-s -w" -o build/linux/harbrixd-amd64  ./cmd/harbrixd

	GOOS=linux GOARCH=arm64   $(GO) build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o build/linux/harbrix-arm64   ./cmd/harbrix
	GOOS=linux GOARCH=arm64   $(GO) build -trimpath -ldflags="-s -w" -o build/linux/harbrixd-arm64  ./cmd/harbrixd

	GOOS=linux GOARCH=arm     $(GO) build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o build/linux/harbrix-arm     ./cmd/harbrix
	GOOS=linux GOARCH=arm     $(GO) build -trimpath -ldflags="-s -w" -o build/linux/harbrixd-arm    ./cmd/harbrixd

	GOOS=linux GOARCH=386     $(GO) build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o build/linux/harbrix-386     ./cmd/harbrix
	GOOS=linux GOARCH=386     $(GO) build -trimpath -ldflags="-s -w" -o build/linux/harbrixd-386    ./cmd/harbrixd


install: build
	$(INSTALL) -d $(BINDIR)
	$(INSTALL) -m 755 build/harbrix $(BINDIR)/
	$(INSTALL) -m 755 build/harbrixd $(BINDIR)/

uninstall:
	rm -f $(BINDIR)/harbrix $(BINDIR)/harbrixd

clean:
	rm -rf build/

run-daemon: build
	sudo ./build/harbrixd

run-cli: build
	./build/harbrix $(ARGS)

list: build
	./build/harbrix list

status: build
	@if [ -z "$(SERVICE)" ]; then echo "Usage: make status SERVICE=<name>"; exit 1; fi
	./build/harbrix status $(SERVICE)

start: build
	@if [ -z "$(SERVICE)" ]; then echo "Usage: make start SERVICE=<name>"; exit 1; fi
	./build/harbrix start $(SERVICE)

stop: build
	@if [ -z "$(SERVICE)" ]; then echo "Usage: make stop SERVICE=<name>"; exit 1; fi
	./build/harbrix stop $(SERVICE)