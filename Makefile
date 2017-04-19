# ref. http://postd.cc/auto-documented-makefile/

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## compile app
	go build cmd/gumk/gumk.go

install: ## Install packages for dependencies
	glide install
