package server

import (
	"context"
	"log"
	"net/http"
	"time"

	// import http profilling when the server profilling configuration is set
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/around25/products/matching-engine/net"

	"github.com/prometheus/client_golang/prometheus"
)

// Server interface
type Server interface {
	Listen()
}

type server struct {
	config  Config
	ctx     context.Context
	markets map[string]MarketEngine
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
	for key, marketCfg := range config.Markets {
		marketEngineConfig := MarketEngineConfig{
			config:   marketCfg,
			producer: NewProducer(config.Brokers.Producers[marketCfg.Publish.Broker], marketCfg.Publish.Topic),
			consumer: NewConsumer(config.Brokers.Consumers[marketCfg.Listen.Broker], marketCfg.Listen.Topic),
		}
		markets[key] = NewMarketEngine(marketEngineConfig)
	}

	return &server{
		config:  config,
		ctx:     context.Background(),
		markets: markets,
	}
}

// Listen for new events that affect the market and process them
func (srv *server) Listen() {
	// start all markets and listen for incomming events
	for _, market := range srv.markets {
		market.Start(srv.ctx)
	}

	// listen for messages and ditribute them to the correct markets
	go srv.ReceiveMessages()

	srv.StartProfilling(srv.config.Server.Monitoring)
	srv.stopOnSignal()
}

func (srv *server) closeMarkets() {
	for _, market := range srv.markets {
		market.Close()
	}
}

func (srv *server) stopOnSignal() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	sig := <-sigc
	log.Printf("Caught signal %s: Shutting down...", sig)
	srv.closeMarkets()
	time.Sleep(time.Second)
	log.Println("Exiting the trading engine... Have a nice day!")
	os.Exit(0)
}

func (srv *server) ReceiveMessages() {
	for key := range srv.markets {
		// listen for incomming events and fwd them to the correct market
		go func(key string, market MarketEngine) {
			for msg := range market.GetMessageChan() {
				// send the message for processing by the specific market
				market.Process(msg)
			}
			log.Println("Closing consumer processor. Message channel closed")
		}(key, srv.markets[key])
	}
}

// StartProfilling enables the engine to be pulled for metrics by prometheus or golang profilling
func (srv server) StartProfilling(config MonitoringConfig) {
	if config.Enabled {
		go func() {
			log.Printf("Starting profilling server and listening %s:%s", config.Host, config.Port)
			http.Handle("/metrics", prometheus.Handler())
			log.Println(http.ListenAndServe(config.Host+":"+config.Port, nil))
		}()
	}
}

// NewConsumer starts a new consumer based on the config
func NewConsumer(config ConsumerConfig, topic string) net.KafkaConsumer {
	return net.NewKafkaConsumer(config.Hosts, topic, 0)
}

// NewProducer starts a new producer based on the config
func NewProducer(config ProducerConfig, topic string) net.KafkaProducer {
	return net.NewKafkaProducer(config.Hosts, topic)
}
