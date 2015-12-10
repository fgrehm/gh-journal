PATH       := node_modules/.bin:$(PATH)
SHELL      := /usr/bin/env bash
# ------------------------------------------------------------------------------

source     := src
source_js   = $(shell find -L $(source) -type f -name '*.js')
source_css  = $(shell find -L $(source) -type f -name '*.scss')
source_rest = $(shell find -L $(source) ! -name '*.js' ! -name '*.scss' -type f)

output     := ../www
output_css  = $(output)/css/index.css
output_js   = $(output)/js/index.js
output_rest = $(patsubst $(source)/%,$(output)/%,$(source_rest))

output_js_dependencies = node_modules $(source_js)

.PHONY: all
all: build

.PHONY: clean
clean:; rm -rf $(output) browserify-cache.json

.PHONY: build
build: javascript stylesheet rest

.PHONY: build.watch
build.watch:
	$(MAKE) build || true
	watchf -e "write,remove,create" -c "clear" -c "make build" -include ".(js|scss|html)$$" -r

.PHONY: javascript
javascript: $(output_js)

node_modules: package.json
	@npm install
	@touch $@

$(output_js): $(output_js_dependencies)
	@mkdir -p $(dir $@)
	@echo 'Executing browserifyinc...'
	@browserifyinc $(source)/index.js \
		--verbose \
		--outfile $@

.PHONY: stylesheet
stylesheet: $(output_css)

$(output_css): node_modules $(source_css)
	@mkdir -p $(dir $@)
	@echo 'Executing node-sass...'
	@node-sass $(source)/index.scss | postcss \
		--use autoprefixer \
		--autoprefixer.browsers 'last 2 versions' \
		--output $@

.PHONY: rest
rest: $(output_rest)

$(output)/%: $(source)/%
	@mkdir -p $(dir $@)
	cp $< $@
