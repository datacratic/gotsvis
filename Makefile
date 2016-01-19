all: build verify test
verify: lint
test: test-cover test-race
.PHONY: all verify test

fmt:
	@echo -- format source code
	@go fmt ./...
.PHONY: fmt

build: fmt
	@echo -- build all packages
	@go install ./...
.PHONY: build

vet: build
	@echo -- static analysis
	@go vet ./...
.PHONY: vet

lint:
	@echo -- report coding style issues
	@find . -type f -name "*.go" -exec golint {} \;
.PHONY: lint

test-cover:
	@echo -- build and run tests
	@go test -cover -test.short ./...
.PHONY: test-cover

test-race:
	@echo -- rerun all tests with race detector
	@GOMAXPROCS=4 go test -test.short -race ./...
.PHONY: test-race

test-all: 
	@echo -- build and run all tests
	@GOMAXPROCS=4 go test -race ./...
.PHONY:test-all
