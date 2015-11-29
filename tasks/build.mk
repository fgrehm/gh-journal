.PHONY: build
build: bin/gh-journal

bin/gh-journal: $(shell find -L src -type f -name '*.go')
	gb build cmd/gh-journal
