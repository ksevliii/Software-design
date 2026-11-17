DATA_DIR ?= data

.PHONY: run run-strict run-zero test vet fmt check build demo

run: run-strict

run-strict:
	@go run ./cmd/pos -data "$(DATA_DIR)" -mode strict

run-zero:
	@go run ./cmd/pos -data "$(DATA_DIR)" -mode missing_as_zero

test:
	@go test ./...

vet:
	@go vet ./...

fmt:
	@gofmt -w .

check: test vet

build:
	@go build -o pos ./cmd/pos

demo:
	@go run ./cmd/demo