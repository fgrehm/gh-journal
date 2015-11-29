.PHONY: hack
hack:
	docker-compose -f docker-compose.dev.yml run --rm hack

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
