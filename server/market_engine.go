package server

import (
	"log"
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
}

// NewMarketEngine open a new market
func NewMarketEngine(producer net.KafkaProducer, topic string) MarketEngine {
	return &marketEngine{
		producer:     producer,
		publishTopic: topic,
		engine:       trading_engine.NewTradingEngine(),
		inputs:       make(chan *sarama.ConsumerMessage, 10000),
		orders:       make(chan trading_engine.Order, 20000),
		trades:       make(chan []trading_engine.Trade, 20000),
		messages:     make(chan []byte, 20000),
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
		// Monitor: Increment the number of messages that has been received by the market
		messagesQueued.WithLabelValues(mkt.publishTopic).Inc()
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
		messagesQueued.WithLabelValues(mkt.publishTopic).Dec()
		// Monitor: Increment the number of orders that are waiting to be processed
		ordersQueued.WithLabelValues(mkt.publishTopic).Inc()
		mkt.orders <- order
	}
}

// ProcessOrder process each order by the trading engine and forward trades to the trades channel
//
// Message flow is unidirectional from the orders channel to the trades channel
func (mkt *marketEngine) ProcessOrder() {
	for order := range mkt.orders {
		// Process each order and generate trades
		generatedTrades := mkt.engine.Process(order)
		// Monitor: Update order count for monitoring with prometheus
		engineOrderCount.WithLabelValues(mkt.publishTopic).Inc()
		ordersQueued.WithLabelValues(mkt.publishTopic).Dec()
		tradeCount := len(generatedTrades)
		if tradeCount > 0 {
			tradesQueued.WithLabelValues(mkt.publishTopic).Add(float64(tradeCount))
			// send trades for storage
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
		}

		// Monitor: Update the number of trades processed after sending them back to Kafka
		tradeCount := float64(len(completedTrades))
		tradesQueued.WithLabelValues(mkt.publishTopic).Sub(tradeCount)
		engineTradeCount.WithLabelValues(mkt.publishTopic).Add(tradeCount)
	}
}
