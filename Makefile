.PHONY: build install clean
VERSION=`egrep -o '[0-9]+\.[0-9a-z.\-]+' version.go`
GIT_SHA=`git rev-parse --short HEAD || echo`
BUILT=`date +%FT%T`

build:
	@echo "Building oreilly_kindle_book..."
	@mkdir -p bin
	@go build -ldflags "-X main.Built=${BUILT} -X main.GitSHA=${GIT_SHA}" -o bin/oreilly_kindle_book .
	@echo "Building success..."

install:
	@echo "Installing oreilly_kindle_book..."
	@install -c bin/oreilly_kindle_book /usr/local/bin/oreilly_kindle_book
	@echo "Install success to /usr/local/bin/oreilly_kindle_book."
	@oreilly_kindle_book -help

clean:
	@rm -f bin/*
