.PHONY: build doc-publish doc-serve

.DEFAULT_GOAL := help
NAME = ian
GOARCH = amd64
PKG := github.com/thylong/ian
VERSION := $(shell git describe --abbrev=0 --tags)

.PHONY: build
build: ## build Go binaries for targeted platforms
	# Build for MacOS amd64
	GOOS=darwin GOARCH=${GOARCH} go build -i -v -o ${NAME}-darwin-${GOARCH} -ldflags="-X main.version=${VERSION}" ${PKG}
	# Build for Linux amd64
	GOOS=linux GOARCH=${GOARCH} go build -i -v -o ${NAME}-linux-${GOARCH} -ldflags="-X main.version=${VERSION}" ${PKG}
	# Build for Windows amd64
	go get github.com/inconshreveable/mousetrap
	GOOS=windows GOARCH=${GOARCH} go build -i -v -o ${NAME}-windows-${GOARCH}.exe -ldflags="-X main.version=${VERSION}" ${PKG}

.PHONY: test
test: ## run all tests
	go test -cover ./...

.PHONY: doc-publish
doc-publish: ## create live documentation static assets
	hugo -t hugo-theme-learn -s docs-source -d ../docs

.PHONY: doc-serve
doc-serve: ## serve doc site on localhost
	hugo server --buildDrafts -t hugo-theme-learn -s docs-source -w

# Implements this pattern for autodocumenting Makefiles:
# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
#
# Picks up all comments that start with a ## and are at the end of a target definition line.
help:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'