.PHONY: build doc-publish doc-serve

NAME = ian
GOARCH = amd64
PKG := github.com/thylong/ian
VERSION := $(shell git describe --abbrev=0 --tags)

build:
	# Build for MacOS amd64
	GOOS=darwin GOARCH=${GOARCH} go build -i -v -o ${NAME}-darwin-${GOARCH} -ldflags="-X main.version=${VERSION}" ${PKG}
	# Build for Linux amd64
	GOOS=linux GOARCH=${GOARCH} go build -i -v -o ${NAME}-linux-${GOARCH} -ldflags="-X main.version=${VERSION}" ${PKG}
	# Build for Windows amd64
	GOOS=windows GOARCH=${GOARCH} go build -i -v -o ${NAME}-windows-${GOARCH}.exe -ldflags="-X main.version=${VERSION}" ${PKG}

test:
	go test -cover ./...

doc-publish:
	hugo -t hugo-theme-learn -s docs-source -d ../docs

doc-serve:
	hugo server --buildDrafts -t hugo-theme-learn -s docs-source -w
