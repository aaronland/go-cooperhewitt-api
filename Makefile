CWD=$(shell pwd)
GOPATH := $(CWD)

build:	rmdeps deps fmt bin

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/thisisaaronland/go-cooperhewitt-api; then rm -rf src/github.com/thisisaaronland/go-cooperhewitt-api; fi
	mkdir -p src/github.com/thisisaaronland/go-cooperhewitt-api/client
	cp client/*.go src/github.com/thisisaaronland/go-cooperhewitt-api/client/
	mkdir -p src/github.com/thisisaaronland/go-cooperhewitt-api/endpoint
	cp endpoint/*.go src/github.com/thisisaaronland/go-cooperhewitt-api/endpoint/
	mkdir -p src/github.com/thisisaaronland/go-cooperhewitt-api/response
	cp response/*.go src/github.com/thisisaaronland/go-cooperhewitt-api/response/
	mkdir -p src/github.com/thisisaaronland/go-cooperhewitt-api/shoebox
	cp shoebox/*.go src/github.com/thisisaaronland/go-cooperhewitt-api/shoebox/
	mkdir -p src/github.com/thisisaaronland/go-cooperhewitt-api/template
	cp template/*.go src/github.com/thisisaaronland/go-cooperhewitt-api/template/
	mkdir -p src/github.com/thisisaaronland/go-cooperhewitt-api/util
	cp util/*.go src/github.com/thisisaaronland/go-cooperhewitt-api/util/
	cp api.go src/github.com/thisisaaronland/go-cooperhewitt-api/
	if test ! -d src; then mkdir src; fi
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/briandowns/spinner"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/pretty"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-crawl"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt api.go
	go fmt client/*.go
	go fmt cmd/*.go
	go fmt endpoint/*.go
	go fmt response/*.go
	go fmt shoebox/*.go
	go fmt template/*.go
	go fmt util/*.go

bin:	self
	@GOPATH=$(shell pwd) go build -o bin/ch-api cmd/ch-api.go
	@GOPATH=$(shell pwd) go build -o bin/ch-shoebox-archive cmd/ch-shoebox-archive.go
	@GOPATH=$(shell pwd) go build -o bin/ch-shoebox-render cmd/ch-shoebox-render.go
