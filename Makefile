all: install install_test

install:
	@echo "Installing dependencies"
	go get github.com/golang/dep/cmd/dep
	dep ensure

install_test:
	go get github.com/alecthomas/gometalinter
	go install .
	gometalinter --install

build:
	@echo "Building production binary"
	CGO_ENABLED=0 go build -o ./bin/dailyteedeals -ldflags="-s -w" main.go

lint:
	gometalinter --fast

test:
	echo "" > coverage.txt
	for d in $$(go list ./... | grep -v vendor); do \
		go test -coverprofile=profile.out -covermode=atomic $$d; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out; \
		fi \
	done