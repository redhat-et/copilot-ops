##@ General 

IMG ?= quay.io/copilot-ops/copilot-ops
IMG_TAG ?= latest

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

build: ## Build the copilot-ops file.
	@ echo ▶️ go build
	go build
	@ echo ✅ go build
	@ echo ./copilot-ops -h "# run me!"
.PHONY: build

.PHONY: image
image: ## Build a container image of the current project.
	docker build -t ${IMG}:${IMG_TAG} .
	@ echo ✅ docker build

.PHONY: publish
publish: image ## Publish the container image into Quay.io.
	docker push ${IMG}:${IMG_TAG}
	@ echo ✅ docker push

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

##@ Download utilities

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