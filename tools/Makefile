UNAME_S := $(shell uname -s)

ifeq ($(UNAME_S),Darwin)
    OS_TYPE = mac
else ifeq ($(UNAME_S),Linux)
    OS_TYPE = linux
else
    OS_TYPE = unknown
endif

.PHONY: dep
dep:
	go install github.com/boyter/scc/v3@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
if [ "$(OS_TYPE)" = "mac" ]; then \
	brew install antlr
else
	sudo apt-get install -y antlr4
fi
