SHELL := /bin/bash

SUB_MODULE := persistence/dynamodbpersistence

.PHONY: version
version:
	@if [ -z "$(v)" ]; then \
	  echo "Usage: make version v=vMAJOR.MINOR.PATCH"; \
	  exit 1; \
	fi

	@if ! echo "$(v)" | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' > /dev/null; then \
	  echo "Error: Version must be of the form vMAJOR.MINOR.PATCH (e.g. v1.2.3)"; \
	  exit 1; \
	fi

	@echo "==> Checking existing tags for version $(v) in submodule '$(SUB_MODULE)'..."
	@if module_tag="$(SUB_MODULE)/$(v)" && git rev-parse --verify --quiet "$$module_tag" >/dev/null; then \
	  echo "Error: Tag '$$module_tag' already exists"; \
	  exit 1; \
	fi

	@echo "==> Creating new tag..."
	@if module_tag="$(SUB_MODULE)/$(v)"; then \
	  echo "git tag -a \"$$module_tag\" -m \"Release $$module_tag\""; \
	  git tag -a "$$module_tag" -m "Release $$module_tag"; \
	fi

	@echo "==> Tagged $(SUB_MODULE)/$(v)"
	@echo "Don't forget to do: git push --follow-tags"