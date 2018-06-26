package server

import (
	"log"

	"github.com/spf13/viper"
)

// ConsumerConfig structure
type ConsumerConfig struct {
	Name   string
	Hosts  []string
	Topics []string
}

// ProducerConfig structure
type ProducerConfig struct {
	Hosts  []string
	Topics []string
}

// BrokersConfig structure
type BrokersConfig struct {
	Consumers map[string]ConsumerConfig
	Producers map[string]ProducerConfig
}

// TopicConfig structure
type TopicConfig struct {
	Broker string
	Topic  string
}

// MarketConfig structure
type MarketConfig struct {
	Base  string
	Quote string

	QuoteIncrements float64 `mapstructure:"quote_increments"`
	BaseMin         float64 `mapstructure:"base_min"`
	BaseMax         float64 `mapstructure:"base_max"`

	Listen  TopicConfig
	Publish TopicConfig
}

// Config structure
type Config struct {
	Markets map[string]MarketConfig
	Brokers BrokersConfig
}

// LoadConfig Load server configuration from the yaml file
func LoadConfig(viperConf *viper.Viper) Config {
	var config Config

	err := viperConf.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return config
}
