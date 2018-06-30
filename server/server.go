package server

import (
	"log"
	"net/http"
	// import http profilling when the server profilling configuration is set
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trading_engine/net"

	"github.com/prometheus/client_golang/prometheus"
)

// Server interface
type Server interface {
	Listen()
}

type server struct {
	config       Config
	consumers    map[string]net.KafkaConsumer
	producers    map[string]net.KafkaProducer
	markets      map[string]MarketEngine
	topic2market map[string]string
}

var (
	engineOrderCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "engine_order_count",
		Help: "Trading engine order count",
	}, []string{
		// Which market are the orders from?
		"market",
	})
	engineTradeCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "engine_trade_count",
		Help: "Trading engine trade count",
	}, []string{
		// Which market are the trades from?
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
	tradesQueued = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "engine_trade_queue_count",
		Help: "Number of trades waiting to be processed.",
	}, []string{
		// Which market are the trades from?
		"market",
	})
)

func init() {
	prometheus.MustRegister(engineOrderCount)
	prometheus.MustRegister(engineTradeCount)
	prometheus.MustRegister(messagesQueued)
	prometheus.MustRegister(ordersQueued)
	prometheus.MustRegister(tradesQueued)
}

// NewServer constructor
func NewServer(config Config) Server {
	// init producers
	producers := make(map[string]net.KafkaProducer, len(config.Brokers.Producers))
	for key, brokerCfg := range config.Brokers.Producers {
		producers[key] = NewProducer(brokerCfg)
	}

	// init consumers
	consumers := make(map[string]net.KafkaConsumer, len(config.Brokers.Consumers))
	for key, brokerCfg := range config.Brokers.Consumers {
		consumers[key] = NewConsumer(brokerCfg)
	}

	// start markets
	topic2market := make(map[string]string, len(config.Markets))
	markets := make(map[string]MarketEngine)
	for key, marketCfg := range config.Markets {
		markets[key] = NewMarketEngine(producers[marketCfg.Publish.Broker], marketCfg.Publish.Topic)
		topic2market[marketCfg.Listen.Topic] = key
	}

	return &server{
		config:       config,
		producers:    producers,
		consumers:    consumers,
		markets:      markets,
		topic2market: topic2market,
	}
}

// Listen for new events that affect the market and process them
func (srv *server) Listen() {
	// start all producers
	for _, producer := range srv.producers {
		producer.Start()
	}

	// start all markets
	for _, market := range srv.markets {
		market.Start()
	}

	// start all consumers
	for _, consumer := range srv.consumers {
		consumer.Start()
	}

	go srv.ReceiveMessages()

	srv.StartProfilling(srv.config.Server.Profilling)
	srv.stopOnSignal()
}

func (srv *server) closeConsumers() {
	for _, consumer := range srv.consumers {
		consumer.Close()
	}
}
func (srv *server) closeProducers() {
	for _, producer := range srv.producers {
		producer.Close()
	}
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

	log.Printf("Caught signal %s: Shutting down in 3 seconds...", sig)
	log.Println("Closing all consumers...")
	srv.closeConsumers()
	time.Sleep(time.Second)
	log.Println("Closing all producers...")
	srv.closeProducers()
	time.Sleep(2 * time.Second)
	srv.closeMarkets()
	log.Println("Exiting the trading engine... Have a nice day!")
	os.Exit(0)
}

func (srv *server) ReceiveMessages() {
	for key := range srv.consumers {
		go func(key string, consumer net.KafkaConsumer) {
			for msg := range consumer.GetMessageChan() {
				market := srv.topic2market[msg.Topic]
				consumer.MarkOffset(msg, "")
				srv.markets[market].Process(msg)
			}
		}(key, srv.consumers[key])
	}
}

func (srv server) StartProfilling(config ProfillingConfig) {
	if config.Enabled {
		go func() {
			http.Handle("/metrics", prometheus.Handler())
			log.Println(http.ListenAndServe(config.Host+":"+config.Port, nil))
		}()
	}
}

// NewConsumer starts a new consumer based on the config
func NewConsumer(config ConsumerConfig) net.KafkaConsumer {
	return net.NewKafkaPartitionConsumer(config.Name, config.Hosts, config.Topics)
}

// NewProducer starts a new producer based on the config
func NewProducer(config ProducerConfig) net.KafkaProducer {
	return net.NewKafkaAsyncProducer(config.Hosts)
}
