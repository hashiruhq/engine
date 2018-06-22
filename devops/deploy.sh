#!/bin/bash 

mkdir /srv/go
mkdir /srv/go/src
mkdir /srv/go/bin
mkdir /srv/go/pkg

sudo add-apt-repository ppa:gophers/archive -y
sudo apt-get update -y
sudo apt-get install -y golang-1.10-go unzip

export PATH=/usr/lib/go-1.10/bin:$PATH
export GOPATH=/srv/go
export KAFKA_BROKER=46.101.168.89:9092
export KAFKA_ORDER_TOPIC=trading.order.btc.eth
export KAFKA_ORDER_CONSUMER=trading_engine_btc_eth
export KAFKA_TRADE_TOPIC=trading.trade.btc.eth

# update the new ip in the docker-compose file

# fix the path to the kafka server in configs

cd /srv/go/src
wget -O trading_engine.zip https://data.cloud.around25.net/s/pYEJhmJHNR2p1MS/download
unzip trading_engine.zip
go get -v -d ./...
go build ./...
go install ./...

go test -benchmem -timeout 20s -run=^$ trading_engine/trading_engine -bench ^BenchmarkWithRandomData$


mkdir /srv/trading_engine/linux_amd64
cd /srv/trading_engine/linux_amd64
wget -O trading_engine https://data.cloud.around25.net/s/xTDVDaLzPesymmV/download
chmod +x trading_engine

# create systemd file under: /etc/systemd/system/trading_engine.service
[Unit]
Description=Trading Engine

[Service]
Environment='CATALINA_OPTS=-Dappserver.home=/var/lib/archiva/apache-tomcat-current -Dappserver.base=/var/lib/archiva/apache-tomcat-current'
ExecStart=/srv/trading_engine/linux_amd64/trading_engine

[Install]
WantedBy=multi-user.target

# restart the config
systemctl daemon-reload
systemctl enable trading_engine.service
systemctl start trading_engine.service