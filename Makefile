dirs = ./interfaces/... ./commands/... ./services/... ./bot/... .

all: mumbledj

mumbledj: ## Default action. Builds MumbleDJ.
	@env GO15VENDOREXPERIMENT="1" go build .

test: ## Runs unit tests for MumbleDJ.
	@env GO15VENDOREXPERIMENT="1" go test $(dirs)

clean: ## Removes compiled MumbleDJ binaries.
	@rm -f mumbledj*

install: ## Copies MumbleDJ binary to /usr/local/bin for easy execution.
	@cp -f mumbledj* /usr/local/bin/mumbledj

dist: ## Performs cross-platform builds via gox for multiple Linux platforms.
	@go get -u github.com/mitchellh/gox
	@gox -cgo -osarch="linux/amd64 linux/386"

bindata: ## Regenerates bindata.go with an updated configuration file.
	@go get -u github.com/jteeuwen/go-bindata/...
	@go-bindata config.yaml

help: ## Shows this helptext.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
