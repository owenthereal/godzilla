all: build

build: build-godzilla build-godzilla-parser

build-godzilla:
	go build -o bin/godzilla cmd/godzilla/*

build-godzilla-parser:
	cd parser && npm install && npm run build

test:
	go test $$(go list ./... | grep -v "vendor\|ftest")

ftest: build
	go test -v ./ftest/...
