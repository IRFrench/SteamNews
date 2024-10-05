DEFAULT: help
.PHONY: build

build: ## Build a binary of the service
	CGO_ENABLED=0 go build -o build/news cmd/main.go

run: build ## Run the service
	ETC=config.yaml ./build/news -v -q

help: ## Show commands of the makefile (and any included files)
	@awk 'BEGIN {FS = ":.*?## "}; /^[0-9a-zA-Z_.-]+:.*?## .*/ {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Docker section #

docker.build: ## Build the docker container
	docker build -f dockerfile -t steamnews .

docker.run: ## Run the docker container
	docker run \
	-e ETC=config.yaml \
	-v ./config.yaml:/etc/config.yaml \
	steamnews