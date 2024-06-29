VERSION ?= $(shell git describe --tags --always --dirty)

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	@echo "Running golangci-lint on all sources..."
	@golangci-lint run

.PHONY: fmt
fmt:
	@echo "Running gofmt on all sources..."
	@gofmt -s -l -w .

.PHONY: fmtcheck
fmtcheck:
	@bash -c "diff -u <(echo -n) <(gofmt -d .)"
