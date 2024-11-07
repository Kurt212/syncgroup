BIN_DIR = $(PWD)/bin

GOIMPORTS_BIN = $(BIN_DIR)/goimports
GOLANGCI_BIN = $(BIN_DIR)/golangci-lint

.PHONY: all
all: tidy fmt lint test

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test -v -race -cover -coverprofile=coverage.out -timeout 10s -cpu 1,2,4,8 ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...
	$(GOIMPORTS_BIN) -w .

.PHONY: lint
lint: vet
	$(GOLANGCI_BIN) run -c golangci.yaml ./...

.PHONY: ci
ci: deps lint test

.PHONY: deps
deps:
	go mod download
	mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install golang.org/x/tools/cmd/goimports@v0.26.0
	GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
