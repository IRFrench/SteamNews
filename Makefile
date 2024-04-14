DEFAULT: help
.PHONY: build

build: ## Build a binary of the service
	CGO_ENABLED=0 go build -o build/news cmd/main.go

run: build ## Run the service
	ENV=config.yaml ./build/news -v

help: ## Show commands of the makefile (and any included files)
	@awk 'BEGIN {FS = ":.*?## "}; /^[0-9a-zA-Z_.-]+:.*?## .*/ {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)