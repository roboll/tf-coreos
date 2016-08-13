ORG     := $(shell basename $(realpath ..))
NAME    := $(shell basename $(PWD))

build:
	go build .
.PHONY: build

check:
	go vet $(shell go list ./... | grep -v /vendor/)
.PHONY: check

test:
	go test -v $(shell go list ./... | grep -v /vendor/) -cover -race -p=1
.PHONY: test

tools:
	go get -u github.com/roboll/ghr github.com/mitchellh/gox
.PHONY: tools

cross:
	@mkdir -p dist
	gox -os '!freebsd' -arch '!arm' -output "dist/${NAME}_{{.OS}}_{{.Arch}}"
.PHONY: cross

release: cross
	@ghr -b "`git log ${TAG_PREV}..HEAD --oneline --decorate` [Build Info](${BUILD_URL})" \
		-t ${GITHUB_TOKEN} -u ${ORG} ${TAG} dist
.PHONY: release

TAG      = $(shell git describe --tags --abbrev=0 HEAD)
TAG_PREV = $(shell git describe --tags --abbrev=0 HEAD^)
