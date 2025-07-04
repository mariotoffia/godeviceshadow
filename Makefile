.PHONY: test
test:
	@go test ./... -cover
integration-test:
	@go test -tags=integration ./... -cover

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

# Update references and commit changes
	@$(MAKE) update-refs v=$(v)
	@git add .
	@git commit -m "Update dependencies to $(v)"
	@git push

.PHONY: update-refs
update-refs:
# Update submodule references to root module version
	@if [ -z "$(v)" ]; then \
		echo "Usage: make update-refs v=vX.Y.Z"; \
		exit 1; \
	fi
	@if ! echo "$(v)" | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' > /dev/null; then \
		echo "Error: Version must be of the form vMAJOR.MINOR.PATCH (e.g. v1.2.3)"; \
		exit 1; \
	fi
	$(eval ROOT_MODULE := $(shell head -1 go.mod | awk '{print $$2}'))
	@echo "Updating references to root module $(ROOT_MODULE) to version $(v)"
	@find . -type f -name go.mod ! -path './go.mod' | while read -r modfile; do \
		dir=$$(dirname "$$modfile"); \
		(cd "$$dir" && \
			if go list -m $(ROOT_MODULE) >/dev/null 2>&1; then \
				echo "Updating $$dir/go.mod"; \
				go mod edit -require "$(ROOT_MODULE)@$(v)"; \
				go mod tidy; \
			else \
				echo "No reference to $(ROOT_MODULE) in $$dir/go.mod"; \
			fi); \
	done