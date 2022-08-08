build:
	@ echo ▶️ go build
	go build
	@ echo ✅ go build
	@ echo ./copilot-ops -h "# run me!"
.PHONY: build

image:
	docker build -t quay.io/copilot-ops/copilot-ops .
	@ echo ✅ docker build
.PHONY: image

publish: image
	docker push quay.io/copilot-ops/copilot-ops
	@ echo ✅ docker push
.PHONY: publish

##@ Development

.PHONY: lint
lint: golangci-lint ## Lint source code
	@ echo "▶️ golangci-lint run"
	$(GOLANGCILINT) run ./...
	@ echo "✅ golangci-lint run"

.PHONY: test
test: lint ginkgo ## Run tests.
	@ echo "▶️ ginkgo test"
	$(GINKGO) --coverprofile "cover.out" ./...
	@ echo "✅ ginkgo test"

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
	@ echo "▶️ Downloading golangci-lint"
	curl -sSfL $(GOLANGCI_URL) | sh -s -- -b $(LOCALBIN) $(GOLANGCI_VERSION)
	@ echo "✅ Downloading golangci-lint"

.PHONY: ginkgo
GINKGO := $(LOCALBIN)/ginkgo
ginkgo: $(GINKGO) ## Download ginkgo
$(GINKGO): $(LOCALBIN)
	@ echo "▶️ Downloading ginkgo@v2"
	GOBIN=$(LOCALBIN) go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo
	@ echo "✅ Downloaded ginkgo"


##@ Build Dependencies
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@ echo "▶️ Local binary directory not present, creating..."
	mkdir -p $(LOCALBIN)
	@ echo "✅ Local binary directory created"