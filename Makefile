TOOLS_DIR := tools

.PHONY: test
test:
	@go test ./... -cover

.PHONY: dep
dep:
#	@cd $(TOOLS_DIR) && go mod tidy

	@cd $(TOOLS_DIR) && \
	  go list -m -f '{{if (and (not .Indirect) (not .Main))}}{{.Path}}@{{.Version}}{{end}}' all | \
	  grep -v '^$$' | \
	  xargs -t -n1 go install

.PHONY: stats
stats:
	@scc . --include-ext=go,adoc,ts --exclude-dir=.git,node_modules