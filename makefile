CGO_ENABLED=0
BUILD_LOCATION ?= "."

all: install build
install: install_dependencies install_test_dependencies
build: build_migrate build_dailyteedeals

install_dependencies:
	@echo "Installing dependencies"
	go get github.com/golang/dep/cmd/dep
	dep ensure

install_test_dependencies:
	@echo "Installing test dependencies"
	go get github.com/onsi/ginkgo/ginkgo

build_migrate:
	@echo "Building migrate"
	go build -o $(BUILD_LOCATION)/migrate -ldflags="-s -w" migrations/*.go

build_dailyteedeals:
	@echo "Building dailyteedeals"
	go build -o $(BUILD_LOCATION)/dailyteedeals -ldflags="-s -w"

test:
	ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress