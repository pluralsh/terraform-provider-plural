ifndef GOPATH
$(error $$GOPATH environment variable not set)
endif

ifeq (,$(findstring $(GOPATH)/bin,$(PATH)))
$(error $$GOPATH/bin directory is not in your $$PATH)
endif

##@ General

.PHONY: help
help: ## show help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: show-dependency-updates
show-dependency-updates: ## show possible dependency updates
	go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}} {{.Version}} -> {{.Update.Version}}{{end}}' -m all

.PHONY: update-dependencies
update-dependencies: ## update dependencies
	go get -u ./...
	go mod tidy

.PHONY: install-tools
install-tools: ## install required tools
	@cat internal/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI {} go install {}

##@ Build

.PHONY: build
build: ## build
	go build -o terraform-provider-plural

##@ Codegen

.PHONY: generate-docs
generate-docs: install-tools ## generate docs
	tfplugindocs generate

.PHONY: validate-docs
validate-docs: install-tools ## validate generated docs
	tfplugindocs validate

##@ Tests

.PHONY: test
test: ## run tests
	go test ./... -v

.PHONY: lint
lint: install-tools ## run linters
	golangci-lint run ./...

.PHONY: fix
fix: install-tools ## fix issues found by linters
	golangci-lint run --fix ./...
