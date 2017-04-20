DIST_DIR        = dist
GOX_BINARY_PATH = $(DIST_DIR)/{{.Dir}}-{{.OS}}-{{.Arch}}
GOX_OS          = darwin linux windows
GOX_ARCH        = 386 amd64

# ref. http://postd.cc/auto-documented-makefile/

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## compile app
	go build cmd/gumk/gumk.go

install: ## Install packages for dependencies
	glide install

release: ## Release binaries on GitHub by the specified tag
ifeq ($(CIRCLE_TAG),)
	$(warning No CIRCLE_TAG environmental variable)
else
	$(call cross-compile)
	: Releasing binaries on tag: $(CIRCLE_TAG)
	go get github.com/tcnksm/ghr
	@ghr -u delphinus -replace -prerelease -debug $(CIRCLE_TAG) dist/
endif

# $(call cross-compile)
define cross-compile
	rm -fr $(DIST_DIR)
	: cross compile
	go get github.com/mitchellh/gox
	gox -output '$(GOX_BINARY_PATH)' -os '$(GOX_OS)' -arch '$(GOX_ARCH)'
	: archive each binary
	for i in dist/*; \
	do \
		j=$$(echo $$i | sed -e 's/_[^.]*//'); \
		mv $$i $$j; \
		zip -j $${i%.*} $$j; \
		rm $$j;
	done
endef
