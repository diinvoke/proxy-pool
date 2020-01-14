.PHONY: build clean test test-race

VERSION=2.0.1
BIN=proxypool
DIR_SRC=./cmd

GO_ENV=CGO_ENABLED=0
GO_FLAGS=-ldflags="-X main.version=$(VERSION) -X 'main.buildTime=`date`' -extldflags -static"
GO=$(GO_ENV) $(shell which go)
GOROOT=$(shell `which go` env GOROOT)
GOPATH=$(shell `which go` env GOPATH)
DOCKER_PUBLISH_TAG=docker.pkg.github.com/mingcheng/proxypool/proxypool:$(VERSION)

build: $(DIR_SRC)/api.go
	@$(GO) build $(GO_FLAGS) -o $(BIN) $(DIR_SRC)

docker_image: clean
	@docker build -f ./Dockerfile -t proxypool:$(VERSION) .

docker_image_publish: clean docker_image
	@docker tag proxypool:$(VERSION)  $(DOCKER_PUBLISH_TAG)
	@docker push $(DOCKER_PUBLISH_TAG)

install: build
	@$(GO) install $(GO_FLAGS) $(DIR_SRC)

test:
	@$(GO) test .

dist: clean
	@goreleaser  --skip-publish --rm-dist --snapshot

release:
	@goreleaser --rm-dist

test-race:
	@$(GO) test -race .

protoc:
	@[ -f $(shell go env GOPATH)/bin/protoc-gen-go ] || go get -u github.com/golang/protobuf/protoc-gen-go

protobuf: protoc
	protoc --go_out=plugins=grpc:. ./protobuf/proxypool.proto

# clean all build result
clean:
	@$(GO) clean ./...
	@rm -f $(BIN)
	@rm -rf ./dist/*
