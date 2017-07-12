all: build

build: build-godzilla build-godzilla-parser

install: build-godzilla build-godzilla-parser
	install bin/godzilla bin/godzillac bin/godzilla-parser $$GOPATH/bin

build-godzilla:
	go build -o bin/godzilla cmd/godzilla/*
	go build -o bin/godzillac cmd/godzillac/*

build-godzilla-parser:
	cd parser && npm install && npm run build

test:
	go test $$(go list ./... | grep -v "vendor\|ftest")

ftest: build
	go test -v ./ftest/...
