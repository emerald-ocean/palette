PKGS := ./cmd/... ./internal/... ./pkg/...
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0

.PHONY: build vet test test-race lint check

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" $(PKGS)

vet:
	go vet $(PKGS)

test:
	go test $(PKGS)

test-race:
	go test -race $(PKGS)

lint:
	$(GOLANGCI_LINT) run $(PKGS)

check:
	go run ./cmd/emerald-gen check
