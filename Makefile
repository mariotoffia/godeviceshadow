.PHONY: test
test:
	@go test ./... -cover

.PHONY: dep
dep:
	@make -C tools dep

.PHONY: stats
stats:
	@scc . --include-ext=go,adoc,ts --exclude-dir=.git,node_modules

.PHONY: lint
lint:
	@golangci-lint run ./...