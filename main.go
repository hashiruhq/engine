package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"trading_engine/config"
	"trading_engine/net"
	"trading_engine/trading_engine"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	engine := trading_engine.NewTradingEngine()

	kafkaBroker := config.Config.Get("KAFKA_BROKER")
	kafkaOrderTopic := config.Config.Get("KAFKA_ORDER_TOPIC")
	kafkaOrderConsumer := config.Config.Get("KAFKA_ORDER_CONSUMER")
	kafkaTradeTopic := config.Config.Get("KAFKA_TRADE_TOPIC")

	log.Printf("Connecting to Kafka Broker: %s", kafkaBroker)
	log.Printf("And listening on order topic: %s", kafkaOrderTopic)
	log.Printf("With the consumer name: %s", kafkaOrderConsumer)
	log.Printf("And publishing trades on this trade topic: %s", kafkaTradeTopic)

	orders := make(chan trading_engine.Order, 10000)
	defer close(orders)

	trades := make(chan []trading_engine.Trade, 10000)
	defer close(trades)

	messages := make(chan []byte, 10000)
	defer close(messages)

	producer := net.NewKafkaAsyncProducer([]string{kafkaBroker}, kafkaTradeTopic)
	producer.Start()
	defer producer.Close()

	consumer := net.NewKafkaPartitionConsumer([]string{kafkaBroker}, []string{kafkaOrderTopic})
	consumer.Start(kafkaOrderConsumer)
	defer consumer.Close()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	// listen for producer errors when publishing trades
	go ReceiveProducerErrors(producer)
	// receive messages from the kafka server
	go ReceiveMessages(consumer, messages)
	// decode the json value for each message received into an Order Structure
	go DecodeMessage(messages, orders)
	// process each order by the trading engine and forward trades to the trades channel
	go ProcessOrder(engine, orders, trades)
	// publish trades to the kafka server
	go PublishTrades(producer, trades)

	closeOnSignal(producer, consumer, sigc)
}

func closeOnSignal(producer net.KafkaProducer, consumer net.KafkaConsumer, sigc chan os.Signal) {
	sig := <-sigc
	log.Printf("Caught signal %s: Shutting down in 3 seconds...", sig)
	log.Println("Closing Kafka consumer client...")
	consumer.Close()
	time.Sleep(time.Second)
	log.Println("Closing the Kafka producer...")
	producer.Close()
	time.Sleep(2 * time.Second)
	log.Println("Exiting the trading engine... Have a nice day!")
	os.Exit(0)
}

// ReceiveProducerErrors listens for error messages trades that have not been published
// and sends them on another queue to processing later.
//
// @todo Optionally we can send them to another service, like an in memory database (Redis) so
// that we don't retry to push then to a broken Kafka instance
func ReceiveProducerErrors(producer net.KafkaProducer) {
	errors := producer.Errors()
	for err := range errors {
		log.Print("Error received from trades producer", err)
	}
}

// ReceiveMessages waits for new messages from the consumer and sends them
// to the messages channel for processing
//
// Message flow is unidirectional from the Kafka Consumer to the messages channel
// When the consumer is closed the messages channel can also be closed and we can shutdown the engine
func ReceiveMessages(consumer net.KafkaConsumer, messages chan<- []byte) {
	msgChan := consumer.GetMessageChan()
	for {
		msg, more := <-msgChan
		if !more {
			return
		}
		consumer.MarkOffset(msg, "")
		messages <- msg.Value
	}
}

// DecodeMessage decode the json value for each message received into an Order struct
// before sending it for processing by the trading engine
//
// Message flow is unidirectional from the messages channel to the orders channel
func DecodeMessage(messages <-chan []byte, orders chan<- trading_engine.Order) {
	for msg := range messages {
		var order trading_engine.Order
		order.FromJSON(msg)
		orders <- order
	}
}

// ProcessOrder process each order by the trading engine and forward trades to the trades channel
//
// Message flow is unidirectional from the orders channel to the trades channel
func ProcessOrder(engine *trading_engine.TradingEngine, orders <-chan trading_engine.Order, trades chan<- []trading_engine.Trade) {
	for order := range orders {
		generatedTrades := engine.Process(order)
		if len(generatedTrades) > 0 {
			trades <- generatedTrades
		}
	}
}

// PublishTrades listens for new trades from the trading engine and publishes them to the Kafka server
func PublishTrades(producer net.KafkaProducer, trades <-chan []trading_engine.Trade) {
	for completedTrades := range trades {
		for _, trade := range completedTrades {
			rawTrade, _ := trade.ToJSON() // @todo thread error on encoding json object (low priority)
			producer.Input() <- rawTrade
		}
	}
}
