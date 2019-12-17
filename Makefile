# PROJECT_NAME := "conditional_engine"
# PKG := "around25.com/exchange/$(PROJECT_NAME)"
# PKG_LIST := $(shell go list ./... | grep -v /vendor/)

# by default execute build and install
all: build install

# build the application to check for any compilation errors
build:
	# gofmt -w ./
	# go vet
	go build ./...

# install all dependencies used by the application
deps:
	go get -v -d ./...
	go get -u golang.org/x/lint/golint
	go get github.com/smartystreets/goconvey
	go get github.com/securego/gosec/cmd/gosec

# install the application in the Go bin/ folder
install:
	go install ./...

check:
	gosec ./...
	golint ./...

test:
	go test -v ./...

test-watch:
	goconvey -port=8081 -cover=true .

coverage-test:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

# install the application for all architectures targeted
install-all:
	GOOS=linux GOARCH=amd64 go install
	GOOS=darwin GOARCH=amd64 go install
	# GOOS=windows GOARCH=amd64 go install
	# GOOS=windows GOARCH=386 go install