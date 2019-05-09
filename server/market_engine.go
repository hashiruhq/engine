package server

import (
	"context"
	"log"
	"time"

	"gitlab.com/around25/products/matching-engine/engine"
	"gitlab.com/around25/products/matching-engine/net"

	"github.com/segmentio/kafka-go"
)

// MarketEngine defines how we can communicate to the trading engine for a specific market
type MarketEngine interface {
	Start(context.Context)
	Close()
	GetMessageChan() <-chan kafka.Message
	LoadMarketFromBackup() error
	Process(kafka.Message)
}

// marketEngine structure
type marketEngine struct {
	name     string
	engine   engine.TradingEngine
	inputs   chan kafka.Message
	messages chan engine.Event
	orders   chan engine.Event
	trades   chan engine.Event
	backup   chan bool
	stats    chan bool
	producer net.KafkaProducer
	consumer net.KafkaConsumer
	config   MarketEngineConfig
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
		name:     config.config.MarketID,
		engine:   engine.NewTradingEngine(config.config.MarketID, config.config.PricePrecision, config.config.VolumePrecision),
		backup:   make(chan bool),
		orders:   make(chan engine.Event, 20000),
		trades:   make(chan engine.Event, 20000),
		messages: make(chan engine.Event, 20000),
	}
}

func (mkt *marketEngine) GetMessageChan() <-chan kafka.Message {
	return mkt.consumer.GetMessageChan()
}

// Start the engine for this market
func (mkt *marketEngine) Start(ctx context.Context) {
	// load last market snapshot from the backup files and update offset for the trading engine consumer
	mkt.LoadMarketFromBackup()
	mkt.producer.Start()
	mkt.consumer.Start(ctx)
	// decode the binary value for each message received into an Order Structure
	go mkt.DecodeMessage()
	// process each order by the trading engine and forward trades to the trades channel
	go mkt.ProcessOrder()
	// publish trades to the kafka server
	go mkt.PublishTrades()
	// start the backup scheduler
	go mkt.ScheduleBackup()
}

// Process a new message from the consumer
func (mkt *marketEngine) Process(msg kafka.Message) {
	// Monitor: Increment the number of messages that has been received by the market
	messagesQueued.WithLabelValues(mkt.name).Inc()
	mkt.messages <- engine.NewEvent(msg)
}

// Close the market by closing all communication channels
func (mkt *marketEngine) Close() {
	close(mkt.messages)
	close(mkt.orders)
	close(mkt.trades)
}

// ScheduleBackup sets up an interval at which to automatically back up the market on Kafka
func (mkt *marketEngine) ScheduleBackup() {
	if mkt.config.config.Backup.Interval == 0 {
		log.Printf("[warn] Backup disabled for '%s' market. Please set backup interval to a positive value.", mkt.name)
		return
	}
	for {
		time.Sleep(time.Duration(mkt.config.config.Backup.Interval) * time.Minute)
		mkt.backup <- true
	}
}

// DecodeMessage decode the binary value for each message received into an Order struct
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
				log.Printf("[Backup] [%s] Snapshot created\n", mkt.name)
			}
		case event, more := <-mkt.orders:
			if !more {
				log.Printf("[info] [market:%s] Closing order processor", mkt.name)
				return
			}
			log.Printf(
				"[%s:%d][%s] %s:%s %d@%d\n",
				event.Order.EventType,
				event.Order.ID,
				event.Order.Market,
				event.Order.Side,
				event.Order.Type,
				event.Order.Amount,
				event.Order.Price,
			)
			trades := make([]engine.Trade, 0, 5)
			// Process each order and generate trades
			mkt.engine.ProcessEvent(event.Order, &trades)
			event.SetTrades(trades)
			lastTopic = event.Msg.Topic
			lastPartition = int32(event.Msg.Partition)
			lastOffset = event.Msg.Offset
			// Monitor: Update order count for monitoring with prometheus
			engineOrderCount.WithLabelValues(mkt.name).Inc()
			ordersQueued.WithLabelValues(mkt.name).Dec()
			tradesQueued.WithLabelValues(mkt.name).Add(float64(len(event.Trades)))
			log.Printf(
				"-> [%s:%d][%s]: generated %d trades\n",
				event.Order.EventType,
				event.Order.ID,
				event.Order.Market,
				len(event.Trades),
			)
			for _, trade := range event.Trades {
				log.Printf("-> [trade] %s ask:%d bid:%d %d@%d\n", trade.MakerSide, trade.AskID, trade.BidID, trade.Amount, trade.Price)
			}
			// send trades for storage
			mkt.trades <- event
		}
	}
}

// PublishTrades listens for new trades from the trading engine and publishes them to the Kafka server
func (mkt *marketEngine) PublishTrades() {
	for event := range mkt.trades {
		trades := make([]kafka.Message, len(event.Trades))
		for index, trade := range event.Trades {
			rawTrade, _ := trade.ToBinary() // @todo add better error handling on encoding
			trades[index] = kafka.Message{
				Value: rawTrade,
			}
		}
		err := mkt.producer.WriteMessages(context.Background(), trades...)
		if err != nil {
			log.Println("[error] [kafka] [market:%s] Unable to publish trades: %v", mkt.name, err)
		}

		// Monitor: Update the number of trades processed after sending them back to Kafka
		tradeCount := float64(len(event.Trades))
		tradesQueued.WithLabelValues(mkt.name).Sub(tradeCount)
		engineTradeCount.WithLabelValues(mkt.name).Add(tradeCount)
	}

	log.Printf("Closing output event/trade processor for '%s' market.", mkt.name)
}
