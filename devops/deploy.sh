#!/bin/bash 
export KAFKA_BROKER=kafka:9092
export KAFKA_ORDER_TOPIC=trading.order.btc.eth
export KAFKA_ORDER_CONSUMER=trading_engine_btc_eth
export KAFKA_TRADE_TOPIC=trading.trade.btc.eth
wget -O trading_engine https://data.cloud.around25.net/s/3OcAs55RPHB4C7b/download
chmod +x trading_engine
