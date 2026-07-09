SHELL := /bin/sh
.SHELLFLAGS := -eu -c
.ONESHELL:
.DELETE_ON_ERROR:

BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
GO := go

.PHONY: build test lint clean pack publish

build: $(BIN_DIR)/cognitiveos-cli

$(BIN_DIR)/cognitiveos-cli:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 $(GO) build -ldflags="-s -w" -o $@ ./cmd/cognitiveos-cli
	@echo "  -> $@"

pack: build
	@VERSION=$$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
	@CPM=/workspace/cpm/build/bin/cpm
	@$${CPM} pack --bin $(BIN_DIR)/cognitiveos-cli --name cognitiveos-cli --version $$VERSION --os linux --arch amd64 --description "CognitiveOS TUI frontend"

publish: pack
	@if [ -z "$${REGISTRY_TOKEN}" ]; then \
		echo "  ERROR: REGISTRY_TOKEN not set"; exit 1; \
	fi
	@VERSION=$$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
	@for cgp in *.cgp; do \
		[ -f "$$cgp" ] || continue; \
		URL="https://github.com/CognitiveOS-Project/cli/releases/download/$$VERSION/$$(basename $$cgp)"; \
		/workspace/cpm/build/bin/cpm publish "$$cgp" --download-url "$$URL"; \
		rm "$$cgp"; \
	done

test:
	$(GO) test ./... -v -count=1

lint:
	shellcheck scripts/build.sh
	$(GO) vet ./...

clean:
	rm -rf $(BUILD_DIR)
