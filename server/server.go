package server

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	// import http profilling when the server profilling configuration is set
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashiruhq/engine/license"
	"github.com/hashiruhq/engine/net"
	"github.com/hashiruhq/engine/version"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const EnvDev = "dev"
const VariantDev = "(Dev)"
const EnvDevMaxOffset int64 = 20000

// Server interface
type Server interface {
	Listen()
}

type server struct {
	config  Config
	ctx     context.Context
	markets map[string]MarketEngine
	checker *license.LicenseChecker
}

var (
	engineOrderCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "engine_order_count",
		Help: "Trading engine order count",
	}, []string{
		// Which market are the orders from?
		"market",
	})
	engineEventCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "engine_trade_count",
		Help: "Trading engine trade count",
	}, []string{
		// Which market are the events from?
		"market",
	})
	messagesQueued = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "engine_message_queue_count",
		Help: "Number of messages from Apache Kafka received and waiting to be processed.",
	}, []string{
		// Which market are the orders from?
		"market",
	})
	ordersQueued = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "engine_order_queue_count",
		Help: "Number of orders waiting to be processed.",
	}, []string{
		// Which market are the orders from?
		"market",
	})
	eventsQueued = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "engine_trade_queue_count",
		Help: "Number of events waiting to be processed.",
	}, []string{
		// Which market are the events from?
		"market",
	})
)

func init() {
	prometheus.MustRegister(engineOrderCount)
	prometheus.MustRegister(engineEventCount)
	prometheus.MustRegister(messagesQueued)
	prometheus.MustRegister(ordersQueued)
	prometheus.MustRegister(eventsQueued)
}

// NewServer constructor
func NewServer(config Config) Server {
	// start markets
	markets := make(map[string]MarketEngine)
	maxOffset := int64(0)
	if config.Environment == EnvDev || version.Variant == VariantDev {
		maxOffset = EnvDevMaxOffset
		log.Info().Str("section", "server").Str("action", "init").Msg("Running in development mode. Limited to max 20k commands per market.")
	}
	for key, marketCfg := range config.Markets {
		marketEngineConfig := MarketEngineConfig{
			config:    marketCfg,
			producer:  NewProducer(config.Kafka.Writer, config.Brokers.Producers[marketCfg.Publish.Broker], config.Kafka.UseTLS, marketCfg.Publish.Topic),
			consumer:  NewConsumer(config.Kafka.Reader, config.Brokers.Consumers[marketCfg.Listen.Broker], config.Kafka.UseTLS, marketCfg.Listen.Topic),
			maxOffset: maxOffset,
		}
		markets[key] = NewMarketEngine(marketEngineConfig)
	}

	return &server{
		config:  config,
		ctx:     context.Background(),
		markets: markets,
		checker: license.NewLicenseChecker(version.ProductID, config.LicenseKey, version.Variant),
	}
}

func (srv *server) checkLicense() {
	cfg := srv.config
	log.Info().Str("section", "server").Str("action", "check_license").Msg("Checking server license")
	uses, err := license.VerifyLicense(version.ProductID, cfg.LicenseKey, version.Variant, false)
	if err != nil {
		log.Info().Str("section", "server").Str("action", "check_license").Msg(err.Error())
		os.Exit(1)
		return
	}
	if uses > version.MaxUses {
		log.Info().Str("section", "server").Str("action", "check_license").Int("max_instances", version.MaxUses).Msg("Max instances exceeded")
		os.Exit(2)
		return
	}
	if len(cfg.Markets) > version.MaxMarkets {
		log.Info().Str("section", "server").Str("action", "check_license").Int("max_markets", version.MaxMarkets).Msg("Max number of markets exceeded")
		os.Exit(3)
		return
	}
	// start periodic license checker
	exitChannel := srv.checker.Start()
	go func(exitChannel chan bool) {
		<-exitChannel
		srv.shutdown(4)
	}(exitChannel)
}

// Listen for new events that affect the market and process them
func (srv *server) Listen() {
	// start prometheus profilling metrics
	go loopProfillingServer(srv.config.Server.Monitoring)
	// check license
	if srv.config.Environment != EnvDev && version.Variant != VariantDev {
		srv.checkLicense()
	}

	// start all markets and listen for incomming events
	for _, market := range srv.markets {
		market.Start(srv.ctx)
	}

	// listen for messages and ditribute them to the correct markets
	go srv.ReceiveMessages()
	srv.stopOnSignal()
}

func (srv *server) closeMarkets() {
	for _, market := range srv.markets {
		market.Close()
	}
}

func (srv *server) stopOnSignal() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGINT)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	<-sigc
	srv.shutdown(0)
}

func (srv *server) shutdown(code int) {
	log.Info().Str("section", "server").Str("action", "terminate").Msg("Received shutdown signal. Starting graceful shutdown...")
	log.Debug().Str("section", "server").Str("action", "terminate").Msg("Markets closing...")
	srv.closeMarkets()
	log.Debug().Str("section", "server").Str("action", "terminate").Msg("Waiting a few seconds to finish")
	time.Sleep(1 * time.Second)
	log.Info().Str("section", "server").Str("action", "terminate").Msg("Exiting")
	os.Exit(code)
}

func (srv *server) ReceiveMessages() {
	for key := range srv.markets {
		// listen for incomming events and fwd them to the correct market
		go loopMarketReceive(key, srv.markets[key])
	}
}

func loopMarketReceive(key string, market MarketEngine) {
	for msg := range market.GetMessageChan() {
		// send the message for processing by the specific market
		market.Process(msg)
	}
	log.Info().
		Str("section", "server").
		Str("action", "terminate").
		Str("goroutine", "server.ReceiveMessages").
		Str("market", key).
		Msg("Closing market consumer channel")
}

func loopProfillingServer(config MonitoringConfig) {
	if !config.Enabled {
		return
	}
	log.Debug().
		Str("section", "server").
		Str("action", "init").
		Str("goroutine", "server.profiling").
		Str("host", config.Host).
		Str("port", config.Port).
		Str("path", "/metrics").
		Msg("Starting profilling server")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(config.Host+":"+config.Port, nil)
	if err != nil {
		log.Error().Err(err).
			Str("section", "server").
			Str("action", "init").
			Str("goroutine", "server.profiling").
			Str("host", config.Host).
			Str("port", config.Port).
			Str("path", "/metrics").
			Msg("Error starting metrics server")
	}
}

// NewConsumer starts a new consumer based on the config
func NewConsumer(rCfg net.KafkaReaderConfig, config ConsumerConfig, useTLS bool, topic string) net.KafkaConsumer {
	return net.NewKafkaConsumer(rCfg, config.Hosts, useTLS, topic, 0)
}

// NewProducer starts a new producer based on the config
func NewProducer(wCfg net.KafkaWriterConfig, config ProducerConfig, useTLS bool, topic string) net.KafkaProducer {
	return net.NewKafkaProducer(wCfg, config.Hosts, useTLS, topic)
}
