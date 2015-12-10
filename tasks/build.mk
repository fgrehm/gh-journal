.PHONY: build
build: bin/gh-journal client.build

bin/gh-journal: $(shell find -L src -type f -name '*.go')
	gb build cmd/gh-journal
