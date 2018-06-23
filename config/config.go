package config

import (
	"os"
	"strconv"
)

// Configuration structure
type Configuration interface {
	Get(string) string
	GetInt(string) int
}

type configuration struct {
	Defaults map[string]string
}

// Config keeps the current configuration values
var Config = configuration{
	Defaults: map[string]string{
		// The Kafka broker server to connect to. Defaults to "kafka:9092"
		"KAFKA_BROKER": "kafka:9092",

		// The type of engine to start
		// This can be "SINGLE" or "MULTIPLE"
		// Using the single engine type make better use of system resources and can ensure an optimal execution speed.
		// Using the multiple engine type allows you to process multiple markets at once at the price of individual market performance.
		// You can choose to deploy in SINGLE mode for the most active markets and in multiple mode for slow moving markets, thus
		// achieving optimal resource allocation
		"ENGINE_TYPE": "SINGLE",

		// Single market topics to connect to.
		// This is the recommended way of deploying in order to ensure optimal performance
		"KAFKA_ORDER_TOPIC":    "trading.order.btc.eth",
		"KAFKA_ORDER_CONSUMER": "trading_engine_btc_eth",
		"KAFKA_TRADE_TOPIC":    "trading.trade.btc.eth",

		// Dispatcher settings multiple market deployment
		// Due to performance considerations multi market deployments are not recommended
		// Try to deploy the trading engine on individual servers in order to reach optimal performance for each market pair
		"MAX_QUEUE":                    "1000",
		"MAX_WORKERS":                  "4",
		"KAFKA_ORDER_TOPIC_PATTERN":    "trading.order.{base}.{market}",
		"KAFKA_ORDER_CONSUMER_PATTERN": "trading_engine_{base}_{market}",
		"KAFKA_TRADE_TOPIC_PATTERN":    "trading.trade.{base}.{market}",
		"MARKETS":                      "btc.rpl,btc.eth,btc.ltc,btc.etn,btc.ada,btc.trn,eth.rpl,eth.rpl,eth.ltc,eth.ada,eth.trn",
	},
}

// Get a configuration value from the config
func (config configuration) Get(key string) string {
	value := os.Getenv(key)
	if value == "" {
		return config.Defaults[key]
	}
	return value
}

// GetInt returns an integer value for the given key
func (config configuration) GetInt(key string) int {
	value, err := strconv.Atoi(config.Get(key))
	if err != nil {
		return 0
	}
	return value
}
