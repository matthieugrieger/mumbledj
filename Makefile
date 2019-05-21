dirs = ./interfaces/... ./commands/... ./services/... ./bot/... .

all: asset mumbledj

mumbledj: ## Default action. Builds MumbleDJ.
	@env go build .

.PHONY: test
test: ## Runs unit tests for MumbleDJ.
	@env go test $(dirs)

.PHONY: coverage
coverage: ## Runs coverage tests for MumbleDJ.
	@env overalls -project=go.reik.pl/mumbledj -covermode=atomic
	@mv overalls.coverprofile coverage.txt

.PHONY: clean
clean: ## Removes compiled MumbleDJ binaries.
	@rm -f mumbledj*

.PHONY: install
install: ## Copies MumbleDJ binary to /usr/local/bin for easy execution.
	@cp -f mumbledj* /usr/local/bin/mumbledj

.PHONY: dist
dist: ## Performs cross-platform builds via gox for multiple Linux platforms.
	@go get -u github.com/mitchellh/gox
	@gox -cgo -osarch="linux/amd64 linux/386"

.PHONY: asset
asset: ## Regenerates assets which will be bundled with binary
	@go get github.com/gobuffalo/packr/v2/packr2
	@packr2

.PHONY: help
help: ## Shows this helptext.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
