build:
	@ echo ▶️\  go build
	go build
	@ echo ✅ go build
	@ echo ./copilot-ops -h "# run me!"
.PHONY: build

test: build
	@ echo ▶️\  go test
	go test ./...
	@ echo ✅ go test
	@ echo ▶️\  go vet
	go vet ./...
	@ echo ✅ go vet
.PHONY: test
