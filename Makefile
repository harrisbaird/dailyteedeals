all: install build

install:
	@echo "Installing dependencies"
	go get github.com/golang/dep/cmd/dep
	go get github.com/alecthomas/gometalinter
	dep ensure
	gometalinter --install

build:
	@echo "Building production binary"
	CGO_ENABLED=0 go build -o ./bin/dailyteedeals -ldflags="-s -w" main.go

lint:
	gometalinter $(go list ./... | grep -v /vendor/)

test:
	ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress