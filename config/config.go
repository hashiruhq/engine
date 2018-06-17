package config

import (
	"os"
	"strconv"
)

// Configuration structure
type Configuration struct {
	Defaults map[string]string
}

// Config keeps the current configuration values
var Config = Configuration{
	Defaults: map[string]string{
		"KAFKA_BROKER":                "kafka:9092",
		"KAFKA_ORDER_TOPIC_PREFIX":    "trading.order.",
		"KAFKA_ORDER_CONSUMER_PREFIX": "trading_engine_",
		"KAFKA_TRADE_TOPIC_PREFIX":    "trading.trade.",

		"MARKETS": "btc.rpl,btc.eth,btc.ltc,btc.etn,btc.ada,btc.trn,eth.rpl,eth.rpl,eth.ltc,eth.ada,eth.trn",

		// Dispatcher settings (unused for main engine)
		"MAX_QUEUE":   "1000",
		"MAX_WORKERS": "4",
	},
}

// Get a configuration value from the config
func (config *Configuration) Get(key string) string {
	value := os.Getenv(key)
	if value == "" {
		return config.Defaults[key]
	}
	return value
}

// GetInt returns an integer value for the given key
func (config *Configuration) GetInt(key string) int {
	value, err := strconv.Atoi(config.Get(key))
	if err != nil {
		return 0
	}
	return value
}
