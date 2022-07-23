.DEFAULT_GOAL := dev

.PHONY: dev
dev: ## dev build
dev: tools generate vet fmt lint mod-tidy

.PHONY: run
run: ## run build
run: run-services
	$(call print-target)
	go run -race .

.PHONY: run-services
run-services: ## run docker compose services
	$(call print-target)
	docker compose up --detach

.PHONY: tools
tools: ## install golang tools
	$(call print-target)
	cd _tools && go install $(shell cd _tools && go list -f '{{ join .Imports " " }}' -tags=tools)

.PHONY: clean
clean: ## go clean build, test and modules caches
	$(call print-target)
	go clean -r -i -cache -testcache -modcache

.PHONY: mod-tidy
mod-tidy: ## go mod tidy
	$(call print-target)
	go mod tidy
	cd _tools && go mod tidy

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## golangci-lint
	$(call print-target)
	golangci-lint run

.PHONY: vet
vet: ## go vet
	$(call print-target)
	go vet ./...

.PHONY: fmt
fmt: ## go fmt
	$(call print-target)
	go fmt ./...

.PHONY: generate
generate: ## go generate
generate: tools
	$(call print-target)
	go generate ./...


define print-target
	@printf "Executing target: \033[36m$@\033[0m\n"
endef