all: install install_test

install:
	@echo "Installing dependencies"
	go get github.com/golang/dep/cmd/dep
	dep ensure

install_test:
	go get github.com/alecthomas/gometalinter
	go get github.com/modocache/gover
	go get github.com/mattn/goveralls
	go get github.com/onsi/ginkgo/ginkgo
	gometalinter --install

build:
	@echo "Building production binary"
	CGO_ENABLED=0 go build -o ./bin/dailyteedeals -ldflags="-s -w" main.go

lint:
	gometalinter $(go list ./... | grep -v /vendor/)

test:
	ginkgo -r --randomizeAllSpecs --randomizeSuites --cover --trace --race --progress