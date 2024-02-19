package server

import (
	"github.com/hashiruhq/engine/net"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// ConsumerConfig structure
type ConsumerConfig struct {
	Hosts []string
}

// ProducerConfig structure
type ProducerConfig struct {
	Hosts []string
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
	MarketID        string `mapstructure:"market_id"`
	PricePrecision  int    `mapstructure:"price_precision"`
	VolumePrecision int    `mapstructure:"volume_precision"`

	QuoteIncrements float64 `mapstructure:"quote_increments"`
	BaseMin         float64 `mapstructure:"base_min"`
	BaseMax         float64 `mapstructure:"base_max"`

	Backup MarketBackupConfig

	Listen  TopicConfig
	Publish TopicConfig
}

// MarketBackupConfig structure
type MarketBackupConfig struct {
	Interval int
	Path     string
}

// ServerConfig structure
type ServerConfig struct {
	Debug      bool
	Monitoring MonitoringConfig
}

// MonitoringConfig structure
type MonitoringConfig struct {
	Enabled bool
	Host    string
	Port    string
}

// Config structure
type Config struct {
	LicenseKey  string `mapstructure:"license_key"`
	Environment string `mapstructure:"env"`
	Markets     map[string]MarketConfig
	Brokers     BrokersConfig
	Server      ServerConfig
	Kafka       net.KafkaConfig
}

// LoadConfig Load server configuration from the yaml file
func LoadConfig(viperConf *viper.Viper) Config {
	var config Config

	err := viperConf.Unmarshal(&config)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse configuration file")
	}
	return config
}
