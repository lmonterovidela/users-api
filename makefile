export GO111MODULE ?= on
PACKAGES = $(shell go list ./...)
PACKAGES_PATH = $(shell go list -f '{{ .Dir }}' ./...)

.PHONY: all
all: check_tools ensure-deps test fmt imports linter

.PHONY: check_tools
check_tools:
	@type "golangci-lint" > /dev/null 2>&1 || echo 'Please install golangci-lint: https://github.com/golangci/golangci-lint#install'

.PHONY: ensure-deps
ensure-deps:
	@echo "=> Syncing dependencies with go mod tidy"
	@go mod tidy

.PHONY: test
test:
	@echo "=> Running tests"
	@docker-compose up -d --remove-orphans
	@go test ./... -covermode=atomic -coverpkg=./... -count=1 -race;\
	exit_code=$$?;\
	docker-compose down -v;\
 	exit $$exit_code

.PHONY: test-cover
test-cover:
	@echo "=> Running tests and generating report"
	@docker-compose up -d --remove-orphans
	@go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out -coverpkg=./... -count=1;\
    exit_code=$$?;\
    docker-compose down -v;\
    if [ $$exit_code -ne 0 ]; then \
      	exit $$exit_code;\
    else \
       go tool cover -html=/tmp/coverage.out; \
    fi

.PHONY: up
up:
	@docker-compose up --remove-orphans -V

.PHONY: fmt
fmt:
	@echo "=> Executing go fmt"
	@go fmt $(PACKAGES)

.PHONY: imports
imports:
	@echo "=> Executing goimports"
	@goimports -w $(PACKAGES_PATH)

# Runs golangci-lint with arguments if provided.
.PHONY: linter
linter:
	@echo "=> Executing golangci-lint$(if $(FLAGS), with flags: $(FLAGS))"
	@golangci-lint run ./... $(FLAGS)