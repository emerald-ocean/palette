GOLANGCI_LINT := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0

.PHONY: build vet test test-race lint check

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" ./...

vet:
	go vet ./...

test:
	go test ./...

test-race:
	go test -race ./...

lint:
	$(GOLANGCI_LINT) run ./...

check:
	go run ./cmd/emerald-gen check
