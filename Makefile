all: build

build: build-godzilla build-godzilla-parser

build-godzilla:
	go build -o bin/godzilla cmd/godzilla/*

build-godzilla-parser:
	cd parser && npm run build
