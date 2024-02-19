# PROJECT_NAME := "conditional_engine"
# PKG := "around25.com/exchange/$(PROJECT_NAME)"
# PKG_LIST := $(shell go list ./... | grep -v /vendor/)

# by default execute build and install
all: build install

build-all: clean-build build-dev build-starter build-premium build-enterprise build-corporate

clean-build:
	rm -Rf ./build

# build the application to check for any compilation errors
build:
	# gofmt -w ./
	# go vet
	go build ./...

build-dev:
	mkdir -p ./build/dev/model
	cp ./.engine.yml.template ./build/dev/.engine.yml
	cp ./Dockerfile.template ./build/dev/Dockerfile
	cp ./prometheus.yml.template ./build/dev/prometheus.yml
	cp ./Readme.md.template ./build/dev/Readme.md
	cp ./docker-compose.yml.template ./build/dev/docker-compose.yml
	cp ./model/event.proto ./build/dev/model/event.proto
	cp ./model/market.proto ./build/dev/model/market.proto
	cp ./model/order.proto ./build/dev/model/order.proto
	cp ./model/trade.proto ./build/dev/model/trade.proto
	cp ./docs/grafana_dashboard.json ./build/dev/grafana_dashboard.json
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix dev \
  	--ldflags "-s -w -X 'github.com/hashiruhq/engine/version.Variant=(Dev)' -X 'github.com/hashiruhq/engine/version.ProductID=rNsKn' -X 'github.com/hashiruhq/engine/version.SMaxUses=0' -X 'github.com/hashiruhq/engine/version.SMaxMarkets=3'" \
  	-o ./build/dev/matching_engine_linux_amd64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix dev \
  	--ldflags "-s -w -X 'github.com/hashiruhq/engine/version.Variant=(Dev)' -X 'github.com/hashiruhq/engine/version.ProductID=rNsKn' -X 'github.com/hashiruhq/engine/version.SMaxUses=0' -X 'github.com/hashiruhq/engine/version.SMaxMarkets=3'" \
  	-o ./build/dev/matching_engine_darwin_amd64
	zip -r ./build/matching_engine_dev.zip ./build/dev

build-starter:
	mkdir -p ./build/starter/model
	cp ./.engine.yml.template ./build/starter/.engine.yml
	cp ./Dockerfile.template ./build/starter/Dockerfile
	cp ./prometheus.yml.template ./build/starter/prometheus.yml
	cp ./Readme.md.template ./build/starter/Readme.md
	cp ./docker-compose.yml.template ./build/starter/docker-compose.yml
	cp ./model/event.proto ./build/starter/model/event.proto
	cp ./model/market.proto ./build/starter/model/market.proto
	cp ./model/order.proto ./build/starter/model/order.proto
	cp ./model/trade.proto ./build/starter/model/trade.proto
	cp ./docs/grafana_dashboard.json ./build/starter/grafana_dashboard.json
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix starter \
  	--ldflags "-s -w -X 'github.com/hashiruhq/engine/version.Variant=(Starter)' -X 'github.com/hashiruhq/engine/version.ProductID=rNsKn' -X 'github.com/hashiruhq/engine/version.SMaxUses=2' -X 'github.com/hashiruhq/engine/version.SMaxMarkets=5'" \
  	-o ./build/starter/matching_engine_linux_amd64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix starter \
  	--ldflags "-s -w -X 'github.com/hashiruhq/engine/version.Variant=(Starter)' -X 'github.com/hashiruhq/engine/version.ProductID=rNsKn' -X 'github.com/hashiruhq/engine/version.SMaxUses=2' -X 'github.com/hashiruhq/engine/version.SMaxMarkets=5'" \
  	-o ./build/starter/matching_engine_darwin_amd64
	zip -r ./build/matching_engine_starter.zip ./build/starter

build-premium:
	mkdir -p ./build/premium/model
	cp ./.engine.yml.template ./build/premium/.engine.yml
	cp ./Dockerfile.template ./build/premium/Dockerfile
	cp ./prometheus.yml.template ./build/premium/prometheus.yml
	cp ./Readme.md.template ./build/premium/Readme.md
	cp ./docker-compose.yml.template ./build/premium/docker-compose.yml
	cp ./model/event.proto ./build/premium/model/event.proto
	cp ./model/market.proto ./build/premium/model/market.proto
	cp ./model/order.proto ./build/premium/model/order.proto
	cp ./model/trade.proto ./build/premium/model/trade.proto
	cp ./docs/grafana_dashboard.json ./build/premium/grafana_dashboard.json
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix premium \
  	--ldflags "-s -w -X 'github.com/hashiruhq/engine/version.Variant=(Premium)' -X 'github.com/hashiruhq/engine/version.ProductID=rNsKn' -X 'github.com/hashiruhq/engine/version.SMaxUses=4' -X 'github.com/hashiruhq/engine/version.SMaxMarkets=25'" \
  	-o ./build/premium/matching_engine_linux_amd64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix premium \
  	--ldflags "-s -w -X 'github.com/hashiruhq/engine/version.Variant=(Premium)' -X 'github.com/hashiruhq/engine/version.ProductID=rNsKn' -X 'github.com/hashiruhq/engine/version.SMaxUses=4' -X 'github.com/hashiruhq/engine/version.SMaxMarkets=25'" \
  	-o ./build/premium/matching_engine_darwin_amd64
	zip -r ./build/matching_engine_premium.zip ./build/premium

build-enterprise:
	mkdir -p ./build/enterprise/model
	cp ./.engine.yml.template ./build/enterprise/.engine.yml
	cp ./Dockerfile.template ./build/enterprise/Dockerfile
	cp ./prometheus.yml.template ./build/enterprise/prometheus.yml
	cp ./Readme.md.template ./build/enterprise/Readme.md
	cp ./docker-compose.yml.template ./build/enterprise/docker-compose.yml
	cp ./model/event.proto ./build/enterprise/model/event.proto
	cp ./model/market.proto ./build/enterprise/model/market.proto
	cp ./model/order.proto ./build/enterprise/model/order.proto
	cp ./model/trade.proto ./build/enterprise/model/trade.proto
	cp ./docs/grafana_dashboard.json ./build/enterprise/grafana_dashboard.json
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix enterprise \
  	--ldflags "-s -w -X 'github.com/hashiruhq/engine/version.Variant=(Enterprise)' -X 'github.com/hashiruhq/engine/version.ProductID=rNsKn' -X 'github.com/hashiruhq/engine/version.SMaxUses=15' -X 'github.com/hashiruhq/engine/version.SMaxMarkets=50'" \
  	-o ./build/enterprise/matching_engine_linux_amd64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix enterprise \
  	--ldflags "-s -w -X 'github.com/hashiruhq/engine/version.Variant=(Enterprise)' -X 'github.com/hashiruhq/engine/version.ProductID=rNsKn' -X 'github.com/hashiruhq/engine/version.SMaxUses=15' -X 'github.com/hashiruhq/engine/version.SMaxMarkets=50'" \
  	-o ./build/enterprise/matching_engine_darwin_amd64
	zip -r ./build/matching_engine_enterprise.zip ./build/enterprise

build-corporate:
	mkdir -p ./build/corporate/model
	cp ./.engine.yml.template ./build/corporate/.engine.yml
	cp ./Dockerfile.template ./build/corporate/Dockerfile
	cp ./prometheus.yml.template ./build/corporate/prometheus.yml
	cp ./Readme.md.template ./build/corporate/Readme.md
	cp ./docker-compose.yml.template ./build/corporate/docker-compose.yml
	cp ./model/event.proto ./build/corporate/model/event.proto
	cp ./model/market.proto ./build/corporate/model/market.proto
	cp ./model/order.proto ./build/corporate/model/order.proto
	cp ./model/trade.proto ./build/corporate/model/trade.proto
	cp ./docs/grafana_dashboard.json ./build/corporate/grafana_dashboard.json
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix corporate \
  	--ldflags "-s -w -X 'github.com/hashiruhq/engine/version.Variant=(Corporate)' -X 'github.com/hashiruhq/engine/version.ProductID=rNsKn' -X 'github.com/hashiruhq/engine/version.SMaxUses=50' -X 'github.com/hashiruhq/engine/version.SMaxMarkets=250'" \
  	-o ./build/corporate/matching_engine_linux_amd64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix corporate \
 	 	--ldflags "-s -w -X 'github.com/hashiruhq/engine/version.Variant=(Corporate)' -X 'github.com/hashiruhq/engine/version.ProductID=rNsKn' -X 'github.com/hashiruhq/engine/version.SMaxUses=50' -X 'github.com/hashiruhq/engine/version.SMaxMarkets=250'" \
  	-o ./build/corporate/matching_engine_darwin_amd64
	zip -r ./build/matching_engine_corporate.zip ./build/corporate

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
