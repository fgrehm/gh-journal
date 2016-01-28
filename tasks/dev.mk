.PHONY: hack
hack:
	docker-compose -f docker-compose.dev.yml run --service-ports --rm hack

.PHONY: run
run: build
	./bin/gh-journal

.PHONY: build.watch
build.watch:
	$(MAKE) build || true
	watchf -e "write,remove,create" -c "clear" -c "make build" -include ".go$$" -r

.PHONY: client.build.watch
client.build.watch:
	@cd client && $(MAKE) build.watch

.PHONY: test
test:
	@echo 'Running tests...'
	gb test ...

.PHONY: test.watch
test.watch:
	$(MAKE) test || true
	watchf -e "write,remove,create" -c "clear" -c "make test" -include ".go$$" -r

.PHONY: watch
watch:
	$(MAKE) build || true
	$(MAKE) test || true
	watchf -e "write,remove,create" -c "clear" -c "make build test" -include ".go$$" -r

.PHONY: fmt
fmt:
	go fmt ./src/...
