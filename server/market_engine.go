package server

import (
	"log"
	"time"
	"trading_engine/net"
	"trading_engine/trading_engine"

	"github.com/Shopify/sarama"
)

// MarketEngine defines how we can communicate to the trading engine for a specific market
type MarketEngine interface {
	Start()
	Close()
	Process(*sarama.ConsumerMessage)
}

// marketEngine structure
type marketEngine struct {
	engine       trading_engine.TradingEngine
	inputs       chan *sarama.ConsumerMessage
	messages     chan []byte
	orders       chan trading_engine.Order
	trades       chan []trading_engine.Trade
	producer     net.KafkaProducer
	publishTopic string

	ordersIn  uint64
	tradesOut uint64
}

// NewMarketEngine open a new market
func NewMarketEngine(producer net.KafkaProducer, topic string) MarketEngine {
	return &marketEngine{
		producer:     producer,
		publishTopic: topic,
		engine:       trading_engine.NewTradingEngine(),
		inputs:       make(chan *sarama.ConsumerMessage, 1000),
		orders:       make(chan trading_engine.Order, 10000),
		trades:       make(chan []trading_engine.Trade, 10000),
		messages:     make(chan []byte, 10000),
	}
}

// Start the engine for this market
func (mkt *marketEngine) Start() {
	// listen for producer errors when publishing trades
	go mkt.ReceiveProducerErrors()
	// receive messages from the kafka server
	go mkt.ReceiveMessages()
	// decode the json value for each message received into an Order Structure
	go mkt.DecodeMessage()
	// process each order by the trading engine and forward trades to the trades channel
	go mkt.ProcessOrder()
	// publish trades to the kafka server
	go mkt.PublishTrades()
	// print stats every 10 seconds
	go mkt.PrintStats()
}

// Process a new message from the consumer
func (mkt *marketEngine) Process(msg *sarama.ConsumerMessage) {
	mkt.inputs <- msg
}

// Close the market by closing all communication channels
func (mkt *marketEngine) Close() {
	close(mkt.inputs)
	close(mkt.messages)
	close(mkt.orders)
	close(mkt.trades)
}

func (mkt *marketEngine) PrintStats() {
	var lastOrderCount uint64
	var lastTradeCount uint64
	for {
		time.Sleep(time.Duration(10) * time.Second)
		orderCount := mkt.ordersIn - lastOrderCount
		tradeCount := mkt.tradesOut - lastTradeCount
		lastOrderCount = mkt.ordersIn
		lastTradeCount = mkt.tradesOut
		log.Printf("Market: %s \nOrders/second: %f\nTrades/second: %f\n", mkt.publishTopic, float64(orderCount)/10, float64(tradeCount)/10)
	}
}

// ReceiveProducerErrors listens for error messages trades that have not been published
// and sends them on another queue to processing later.
//
// @todo Optionally we can send them to another service, like an in memory database (Redis) so
// that we don't retry to push then to a broken Kafka instance
func (mkt *marketEngine) ReceiveProducerErrors() {
	errors := mkt.producer.Errors()
	for err := range errors {
		log.Printf("Error sending trade on '%s' topic: %v", mkt.publishTopic, err)
	}
}

// ReceiveMessages waits for new messages from the consumer and sends them
// to the messages channel for processing
//
// Message flow is unidirectional from the Kafka Consumer to the messages channel
// When the consumer is closed the messages channel can also be closed and we can shutdown the engine
func (mkt *marketEngine) ReceiveMessages() {
	for msg := range mkt.inputs {
		mkt.ordersIn++
		mkt.messages <- msg.Value
	}
}

// DecodeMessage decode the json value for each message received into an Order struct
// before sending it for processing by the trading engine
//
// Message flow is unidirectional from the messages channel to the orders channel
func (mkt *marketEngine) DecodeMessage() {
	for msg := range mkt.messages {
		var order trading_engine.Order
		order.FromJSON(msg)
		mkt.orders <- order
	}
}

// ProcessOrder process each order by the trading engine and forward trades to the trades channel
//
// Message flow is unidirectional from the orders channel to the trades channel
func (mkt *marketEngine) ProcessOrder() {
	for order := range mkt.orders {
		generatedTrades := mkt.engine.Process(order)
		if len(generatedTrades) > 0 {
			mkt.trades <- generatedTrades
		}
	}
}

// PublishTrades listens for new trades from the trading engine and publishes them to the Kafka server
func (mkt *marketEngine) PublishTrades() {
	for completedTrades := range mkt.trades {
		for _, trade := range completedTrades {
			rawTrade, _ := trade.ToJSON() // @todo thread error on encoding json object (low priority)
			mkt.producer.Input() <- &sarama.ProducerMessage{
				Topic: mkt.publishTopic,
				Value: sarama.ByteEncoder(rawTrade),
			}
			mkt.tradesOut++
		}
	}
}
