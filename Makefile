default: build
	@echo "Run make help for info about other make targets"

build:                                   ## Build using host go tools
	./ops/build

test:
	go test -v ./...