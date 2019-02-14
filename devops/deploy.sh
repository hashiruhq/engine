#!/bin/bash 

mkdir -p /root/matching_engine/
cd /root/matching_engine/

# mkdir /srv/go
# mkdir /srv/go/src
# mkdir /srv/go/bin
# mkdir /srv/go/pkg

# sudo add-apt-repository ppa:gophers/archive -y
# sudo apt-get update -y
# sudo apt-get install -y golang-1.10-go unzip

# export PATH=/usr/lib/go-1.10/bin:$PATH
# export GOPATH=/srv/go

# update the new ip in the docker-compose file

# fix the path to the kafka server in configs

# cd /srv/go/src
# wget -O matching_engine.zip https://data.cloud.around25.net/s/pYEJhmJHNR2p1MS/download
# unzip matching_engine.zip
# go get -v -d ./...
# go build ./...
# go install ./...
# 
# go test -benchmem -timeout 20s -run=^$ matching_engine/matching_engine -bench ^BenchmarkWithRandomData$


# mkdir /srv/matching_engine/linux_amd64
# cd /srv/matching_engine/linux_amd64
wget -O matching_engine https://data.cloud.around25.net/s/xTDVDaLzPesymmV/download
chmod +x matching_engine

# create systemd file under: /etc/systemd/system/matching_engine.service
[Unit]
Description=Matching Engine

[Service]
ExecStart=/root/matching_engine/matching_engine

[Install]
WantedBy=multi-user.target

# restart the config
systemctl daemon-reload
systemctl enable matching_engine.service
systemctl start matching_engine.service