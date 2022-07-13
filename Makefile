build:
	@ echo ▶️ go build
	go build
	@ echo ✅ go build
	@ echo ./copilot-ops -h "# run me!"
.PHONY: build

test: build lint
	@ echo ▶️ go test
	go clean -testcache ./...
	go test -v ./...
	@ echo ✅ go test
	@ echo ▶️ go vet
	go vet ./...
	@ echo ✅ go vet
.PHONY: test


##@ Development

.PHONY: lint
lint: golangci-lint ## Lint source code
	@ echo "▶️ golangci-lint run"
	$(GOLANGCILINT) run ./...
	@ echo "✅ golangci-lint run"

# .PHONY: test
# test: lint ginkgo ## Run tests.

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

##@ Download utilities

# for performing lints
.PHONY: golangci-lint
GOLANGCILINT := $(LOCALBIN)/golangci-lint
GOLANGCI_URL := https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh
golangci-lint: $(GOLANGCILINT) ## Download golangci-lint
$(GOLANGCILINT): $(LOCALBIN)
	curl -sSfL $(GOLANGCI_URL) | sh -s -- -b $(LOCALBIN) $(GOLANGCI_VERSION)

