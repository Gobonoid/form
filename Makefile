base_dir := $(abspath $(dir $(mkfile_path)))

.PHONY: generate
generate:
	go install github.com/golang/mock/mockgen@v1.5.0
	@(cd $(mktemp -d) && go install github.com/vektra/mockery/v2@latest)
	go generate ./...

LINTER_EXE := golangci-lint
LINTER := $(GOPATH)/bin/$(LINTER_EXE)
LINT_FLAGS ?= -j 2
LINT_RUN_FLAGS ?= -c .golangci.yml
$(LINTER):
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(GOPATH)/bin v1.41.0
.PHONY: go-lint
go-lint: $(LINTER)
	$(LINTER) $(LINT_FLAGS) run $(LINT_RUN_FLAGS)

.PHONY: lint
lint: go-lint

# test
TEST_FLAGS := -v -cover
.PHONY: test
test: ## Run tests
	sh -c "while ! curl -s http://$(ACCOUNT_API_BASE_URL) > /dev/null; do echo waiting for 3s; sleep 3; done"
	$(BUILDENV) go test $(TEST_FLAGS) ./...

.PHONY: all
all: generate lint test

.PHONY: integration
integration:
	docker-compose up \
    --abort-on-container-exit \
    --exit-code-from test