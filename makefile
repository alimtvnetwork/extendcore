BinariesDirectory = ./bin
WindowsBinariesDirectory = bin
MainDirectory = cmd/main
ConfigDirectory = ./configs
ConfigDirectoryForWindows = configs
GoVersion=v1.17.8
MyApp=cli

.PHONY: clean lint changelog snapshot release
.PHONY: build
.PHONY: deps

target: run-main
all: run-main

EXECUTABLES = git go pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

GoPathX=$(shell echo $GOPATH)
VERSION ?= $(shell git describe --tags `git rev-list --tags --max-count=1`)
BINARY = enumcli
MAIN = main.go

BUILDDIR = build
GITREV = $(shell git rev-parse --short HEAD)
BUILDTIME = $(shell date +'%FT%TZ%z')
GO_BUILDER_VERSION=latest

run-main:
	go run cmd/main/*.go

run-server:
	go run cmd/server/*.go

run-client:
	go run cmd/client/*.go

run-sample:
	go run cmd/sample/*.go

build: export-current-binaries
	go build -o $(shell pwd)/$(BUILDDIR)/$(BINARY) $(shell pwd)/$(MainDirectory)/*.go
	@echo "Build $(BINARY) done."
	@echo "Run \"$(shell pwd)/$(BUILDDIR)/$(BINARY)\" to start $(BINARY)."

run:
	$(shell pwd)/$(BUILDDIR)/$(BINARY)

build-run: build run

run-direct:
	"$(BinariesDirectory)/main"

linux-run:
	cd "$(BinariesDirectory)" && ./main

run-tests:
	cd tests && go test -v
	
cat-ssh:
	cat ~/.ssh/id_rsa.pub

ssh-sample:
	echo "ssh-keygen -t rsa -b 4096 -C 'Your email'"
	
modify-authorized-keys:
	sudo vim ~/.ssh/authorized_keys
	
git-clean-get:
	git reset --hard
	git clean -df
	git status
	git pull

run-docker-build-linux:
	./bin/cli-linux-amd64

run-docker-build-windows:
	./bin/cli-windows-amd64.exe

export-private-key:
	export PRIVATE_KEY=$(cat ~/.ssh/id_rsa | base64)

export-go-version:
	export GO_BUILDER_VERSION=$(GoVersion)

export-current-binaries:
	export PATH=${PATH}:"${PWD}"/bin

gopath-export:
	export GOPATH="/root/go"

all-exports: export-private-key export-go-version export-current-binaries gopath-export

changelog:
	git-chglog $(VERSION) > CHANGELOG.md

debug: all-exports
    # Example : https://t.ly/pJiQ, https://t.ly/JEjg, https://t.ly/5Rre, https://goreleaser.com/install/
	docker run --rm --privileged \
    		-e PRIVATE_KEY=$(PRIVATE_KEY) \
    		-v $PWD:/usr/src/$(MyApp) \
    		-v $GOPATH:/go \
    		-v /var/run/docker.sock:/var/run/docker.sock \
    		-w /usr/src/$(MyApp) \
    		ghcr.io/gythialy/golang-cross:$(GO_BUILDER_VERSION) --snapshot --rm-dist

basic-snapshot:
	cd scripts && sudo chmod +x ./docker-deploy.sh
	cd scripts && sudo sh ./docker-deploy.sh

snapshot: all-exports
    # Example : https://t.ly/pJiQ, https://t.ly/JEjg, https://t.ly/5Rre
	docker run --rm -it --privileged \
		-e PRIVATE_KEY=$(PRIVATE_KEY) \
		-v $(CURDIR):/usr/src/myapp \
		-v $(GOPATH):/go \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-w /usr/src/myapp \
		ghcr.io/gythialy/golang-cross:$(GO_BUILDER_VERSION) bash

release: changelog
	docker run --rm -it --privileged \
		-e GITHUB_TOKEN=$(GITHUB_TOKEN) \
		-e PRIVATE_KEY=$(PRIVATE_KEY) \
		-v $(CURDIR):/golang-cross-example \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(GOPATH)/src:/go/src \
		-w /golang-cross-example \
		ghcr.io/gythialy/golang-cross:$(GO_BUILDER_VERSION) --rm-dist --release-notes=CHANGELOG.md
