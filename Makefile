DIST_DIR        = dist
CMD_DIR         = cmd/gumk
GOX_BINARY_PATH = ../../$(DIST_DIR)/{{.Dir}}-{{.OS}}-{{.Arch}}
GOX_OS          = darwin
GOX_ARCH        = amd64

# ref. http://postd.cc/auto-documented-makefile/

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## compile app
	go build cmd/gumk/gumk.go

install: ## Install packages for dependencies
	glide install

compile: ## Compile binaries
	rm -fr $(DIST_DIR)
	: cross compile
	go get github.com/mitchellh/gox
	cd $(CMD_DIR) && gox -output '$(GOX_BINARY_PATH)' -os '$(GOX_OS)' -arch '$(GOX_ARCH)'
	: archive each binary
	for i in dist/*; \
	do \
		zip -j $${i%.*} $$i; \
		rm $$i; \
	done
