VERSION ?= $(shell git describe --tags --always --dirty)
MYSQL_CONTAINER_NAME ?= sqlb-test-mysql

.PHONY: test
test: test-unit test-func

.PHONY: test-unit
test-unit:
	@go test -v -cover ./...

.PHONY: test-func test-func-reset
test-func: test-func-reset
	@cd internal/testing; \
	MYSQL_HOST="$(shell ./internal/testing/scripts/mysql_get_ip.sh $$MYSQL_CONTAINER_NAME)" go test -v ./...

test-func-reset:
	@./internal/testing/scripts/reset.sh

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
