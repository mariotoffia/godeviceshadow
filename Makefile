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
	@golangci-lint run --enable-all --disable=dup'l,exportloopref ./...

.PHONY: version
version:
# Sets version of the root project - for sub-repositories use their respective 
# make version -v=v{MAJOR}.{MINOR}.{PATCH}

# Usage
	@if [ -z "$(v)" ]; then \
		echo "Usage: make version v=v{MAJOR}.{MINOR}.{PATCH}"; \
		exit 1; \
	fi

# Verify version format: must be v{MAJOR}.{MINOR}.{PATCH}, e.g. v1.2.3
	@if ! echo "$(v)" | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' > /dev/null; then \
		echo "Error: Version must be of the form vMAJOR.MINOR.PATCH (e.g. v1.2.3)"; \
		exit 1; \
	fi

# Check if the tag already exist in this git repository -> fail
	@if git rev-parse "$(v)" >/dev/null 2>&1; then \
		echo "Error: Tag '$(v)' already exists"; \
		exit 1; \
	fi

# Tag the repository && push the tag to the remote repository
	@git tag "$(v)"
	@git push --tags

.PHONY: update-refs
update-refs:
# Iterate sub-repositories where the root module is referenced and update that version
	@for d in $(shell go list -m -f '{{.Dir}}' all); do \
		if [ -f "$$d/go.mod" ]; then \
			sed -i '' -e "s|github.com/$(shell go list -m).*/v[0-9]\+\.[0-9]\+\.[0-9]\+|github.com/$(shell go list -m)/$(v)|g" "$$d/go.mod"; \
		fi; \
	done