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
# Sets the version number on all references and tags the repository for consistent versioning
# even if the repo contains a bunch of sub-modules.
#
# Usage
	@if [ -z "$(v)" ]; then \
		echo "Usage: make version v=v{MAJOR}.{MINOR}.{PATCH} [NO_TAG=true]"; \
		echo "       Where NO_TAG=true will omit the actual git tagging"; \
		exit 1; \
	fi

# Verify version format: must be v{MAJOR}.{MINOR}.{PATCH}, e.g. v1.2.3
	@if ! echo "$(v)" | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' > /dev/null; then \
		echo "Error: Version must be of the form vMAJOR.MINOR.PATCH (e.g. v1.2.3)"; \
		exit 1; \
	fi

# Check if any tags already exist under this version
	@echo "Checking existing tags for version $(v)..."
	@for module in $$(find . -name "go.mod" -exec dirname {} \;); do \
		# Fix the root so it doesn't produce a leading slash
		if [ "$$module" = "." ]; then \
			module_tag="$(v)"; \
		else \
			module_tag="$${module#./}/$(v)"; \
		fi; \
		if git rev-parse "$$module_tag" >/dev/null 2>&1; then \
			echo "Error: Tag '$$module_tag' already exists"; \
			exit 1; \
		fi; \
	done

# Update dependencies in each module
	@echo "Updating internal module dependencies to $(v)..."
	@for module in $$(find . -name "go.mod" -exec dirname {} \;); do \
		echo "Updating dependencies in $$module"; \
		cd $$module; \
		# Filter out modules that share the same import path prefix
		for dep in $$(go list -m all | grep "^$$(go list -m)/" | cut -d' ' -f1); do \
			echo "  - Updating $$dep to $(v)"; \
			go get "$$dep@$(v)"; \
		done; \
		go mod tidy; \
		cd - > /dev/null; \
	done

# Create or display tag commands after updates
	@echo "Tagging modules with version $(v)..."
	@for module in $$(find . -name "go.mod" -exec dirname {} \;); do \
		if [ "$$module" = "." ]; then \
			module_tag="$(v)"; \
		else \
			module_tag="$${module#./}/$(v)"; \
		fi; \
		if [ "$(NO_TAG)" = "true" ]; then \
			echo "Would run: git tag -a '$$module_tag' -m 'Release $$module_tag'"; \
		else \
			echo "Tagging $$module_tag"; \
			git tag -a "$$module_tag" -m "Release $$module_tag"; \
		fi; \
	done

	@echo "Done. If you did not use NO_TAG=true, don't forget to run: git push --follow-tags"