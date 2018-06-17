package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"trading_engine/net"
	"trading_engine/trading_engine"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	engine := trading_engine.NewTradingEngine()

	kafkaBroker := Config.Get("KAFKA_BROKER")
	kafkaOrderTopic := Config.Get("KAFKA_ORDER_TOPIC")
	kafkaOrderConsumer := Config.Get("KAFKA_ORDER_CONSUMER")
	kafkaTradeTopic := Config.Get("KAFKA_TRADE_TOPIC")

	producer := net.NewKafkaProducer([]string{kafkaBroker}, kafkaTradeTopic)
	producer.Start()
	defer producer.Close()

	consumer := net.NewKafkaConsumer([]string{kafkaBroker}, []string{kafkaOrderTopic})
	consumer.Start(kafkaOrderConsumer)
	defer consumer.Close()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	go listenToMessages(producer, consumer, engine)
	closeOnSignal(producer, consumer, sigc)
}

func listenToMessages(producer *net.KafkaProducer, consumer *net.KafkaConsumer, engine *trading_engine.TradingEngine) {
	messageChan := consumer.GetMessageChan()
	for msg := range messageChan {
		orderMsg := net.DecodeOrder(msg.Value)
		order := trading_engine.NewOrder(orderMsg.ID, orderMsg.Price, orderMsg.Amount, orderMsg.Side, orderMsg.Category)
		trades := engine.Process(order)
		consumer.MarkOffset(msg, "")
		for i := range trades {
			trade := trades[i]
			tradeMsg := net.EncodeTrade(trade.MakerOrderID, trade.TakerOrderID, trade.Price, trade.Amount)
			producer.Send(tradeMsg.Encoded)
		}
	}
}

func closeOnSignal(producer *net.KafkaProducer, consumer *net.KafkaConsumer, sigc chan os.Signal) {
	sig := <-sigc
	log.Printf("Caught signal %s: shutting down.", sig)
	producer.Close()
	consumer.Close()
	// @todo: Do some cleanning up here
	os.Exit(0)
}
