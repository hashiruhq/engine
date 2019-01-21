package server

import (
	"log"
	"time"

	"gitlab.com/around25/products/matching-engine/net"
	"gitlab.com/around25/products/matching-engine/trading_engine"

	"github.com/Shopify/sarama"
)

// MarketEngine defines how we can communicate to the trading engine for a specific market
type MarketEngine interface {
	Start()
	Close()
	LoadMarketFromBackup() error
	Process(*sarama.ConsumerMessage)
}

// marketEngine structure
type marketEngine struct {
	name           string
	engine         trading_engine.TradingEngine
	inputs         chan *sarama.ConsumerMessage
	messages       chan trading_engine.Event
	orders         chan trading_engine.Event
	trades         chan trading_engine.Event
	backup         chan bool
	producer       net.KafkaProducer
	backupProducer net.KafkaProducer
	consumer       net.KafkaConsumer
	config         MarketEngineConfig
}

// MarketEngineConfig structure
type MarketEngineConfig struct {
	producer net.KafkaProducer
	consumer net.KafkaConsumer
	config   MarketConfig
}

// NewMarketEngine open a new market
func NewMarketEngine(config MarketEngineConfig) MarketEngine {
	return &marketEngine{
		producer: config.producer,
		consumer: config.consumer,
		config:   config,
		name:     config.config.Base + "_" + config.config.Quote,
		engine:   trading_engine.NewTradingEngine(),
		backup:   make(chan bool),
		orders:   make(chan trading_engine.Event, 20000),
		trades:   make(chan trading_engine.Event, 20000),
		messages: make(chan trading_engine.Event, 20000),
	}
}

// Start the engine for this market
func (mkt *marketEngine) Start() {
	// listen for producer errors when publishing trades
	go mkt.ReceiveProducerErrors()
	// decode the json value for each message received into an Order Structure
	go mkt.DecodeMessage()
	// process each order by the trading engine and forward trades to the trades channel
	go mkt.ProcessOrder()
	// publish trades to the kafka server
	go mkt.PublishTrades()
	// start the backup scheduler
	go mkt.ScheduleBackup()
}

// Process a new message from the consumer
func (mkt *marketEngine) Process(msg *sarama.ConsumerMessage) {
	// Monitor: Increment the number of messages that has been received by the market
	messagesQueued.WithLabelValues(mkt.name).Inc()
	mkt.messages <- trading_engine.NewEvent(msg)
}

// Close the market by closing all communication channels
func (mkt *marketEngine) Close() {
	close(mkt.messages)
	close(mkt.orders)
	close(mkt.trades)
}

// ReceiveProducerErrors listens for error messages trades that have not been published
// and sends them on another queue to processing later.
//
// @todo Optionally we can send them to another service, like an in memory database (Redis) so
// that we don't retry to push then to a broken Kafka instance
func (mkt *marketEngine) ReceiveProducerErrors() {
	errors := mkt.producer.Errors()
	for err := range errors {
		log.Printf("Error sending trade on '%s' market: %v", mkt.name, err)
	}
}

// ScheduleBackup sets up an interval at which to automatically back up the market on Kafka
func (mkt *marketEngine) ScheduleBackup() {
	if mkt.config.config.Backup.Interval == 0 {
		log.Printf("Backup disabled for '%s' market. Please set backup interval to a positive value.", mkt.name)
		return
	}
	for {
		time.Sleep(time.Duration(mkt.config.config.Backup.Interval) * time.Minute)
		mkt.backup <- true
	}
}

// DecodeMessage decode the json value for each message received into an Order struct
// before sending it for processing by the trading engine
//
// Message flow is unidirectional from the messages channel to the orders channel
func (mkt *marketEngine) DecodeMessage() {
	for event := range mkt.messages {
		event.Decode()
		messagesQueued.WithLabelValues(mkt.name).Dec()
		// Monitor: Increment the number of orders that are waiting to be processed
		ordersQueued.WithLabelValues(mkt.name).Inc()
		mkt.orders <- event
	}
	log.Printf("Closing decoder processor for '%s' market.", mkt.name)
}

// ProcessOrder process each order by the trading engine and forward trades to the trades channel
//
// Message flow is unidirectional from the orders channel to the trades channel
func (mkt *marketEngine) ProcessOrder() {
	var lastTopic string
	var lastPartition int32
	var lastOffset int64
	var prevOffset int64
	for {
		select {
		case <-mkt.backup:
			// Generate backup event

			if lastTopic != "" && lastOffset != prevOffset {
				market := mkt.engine.BackupMarket()
				market.Topic = lastTopic
				market.Partition = lastPartition
				market.Offset = lastOffset
				prevOffset = lastOffset
				mkt.BackupMarket(market)
				log.Printf("[Backup] [%s] Snapshot created in /root/backups/%s.json\n", mkt.name, mkt.name)
			} else {
				log.Printf("[Backup] [%s] Skipped. No changes since last backup.\n", mkt.name)
			}
		case event := <-mkt.orders:
			lastTopic = event.Msg.Topic
			lastPartition = event.Msg.Partition
			lastOffset = event.Msg.Offset
			// Process each order and generate trades
			event.SetTrades(mkt.engine.Process(event.Order))
			// Monitor: Update order count for monitoring with prometheus
			engineOrderCount.WithLabelValues(mkt.name).Inc()
			ordersQueued.WithLabelValues(mkt.name).Dec()
			tradesQueued.WithLabelValues(mkt.name).Add(float64(len(event.Trades)))
			// send trades for storage
			mkt.trades <- event
		}
	}
}

// PublishTrades listens for new trades from the trading engine and publishes them to the Kafka server
func (mkt *marketEngine) PublishTrades() {
	for event := range mkt.trades {
		for _, trade := range event.Trades {
			rawTrade, _ := trade.ToJSON() // @todo thread error on encoding json object (low priority)
			mkt.producer.Input() <- &sarama.ProducerMessage{
				Topic: mkt.config.config.Publish.Topic,
				Value: sarama.ByteEncoder(rawTrade),
			}
		}
		// mark the order as completed only after the trades have been sent to the kafka server
		mkt.consumer.MarkOffset(event.Msg, "")

		// Monitor: Update the number of trades processed after sending them back to Kafka
		tradeCount := float64(len(event.Trades))
		tradesQueued.WithLabelValues(mkt.name).Sub(tradeCount)
		engineTradeCount.WithLabelValues(mkt.name).Add(tradeCount)
	}

	log.Printf("Closing output event/trade processor for '%s' market.", mkt.name)
}
