package trading_engine_test

import (
	"fmt"
	"log"
	"testing"
	"time"
	"trading_engine/net"
	"trading_engine/trading_engine"
)

func BenchmarkKafkaConsumerProducer(benchmark *testing.B) {
	engine := trading_engine.NewTradingEngine()

	generateOrdersInKafka(benchmark.N)

	kafkaBroker := "kafka:9092"
	kafkaOrderTopic := "trading.order.btc.eth"
	kafkaTradeTopic := "trading.trade.btc.eth"
	kafkaOrderConsumer := "benchmark_kafka_consumer_producer_test"
	ordersCompleted := 0
	tradesCompleted := 0

	// start the producer service to send new trades to
	producer := net.NewKafkaAsyncProducer([]string{kafkaBroker}, kafkaTradeTopic)
	producer.Start()

	// start the consumer service and listen for new orders
	consumer := net.NewKafkaPartitionConsumer([]string{kafkaBroker}, []string{kafkaOrderTopic})
	consumer.Start(kafkaOrderConsumer)
	defer consumer.Close()

	messages := make(chan []byte, 10000)
	defer close(messages)

	orders := make(chan trading_engine.Order, 10000)
	defer close(orders)

	tradeChan := make(chan []trading_engine.Trade, 10000)

	done := make(chan bool)
	defer close(done)

	finishTrades := make(chan bool)
	defer close(finishTrades)

	// receive messages from the kafka server
	receiveMessages := func(messages chan<- []byte, n int) {
		msgChan := consumer.GetMessageChan()
		for j := 0; j < n; j++ {
			msg := <-msgChan
			consumer.MarkOffset(msg, "")
			messages <- msg.Value
		}
	}

	receiveProducerErrors := func() {
		errors := producer.Errors()
		for err := range errors {
			value, _ := err.Msg.Value.Encode()
			log.Print("Error received from trades producer", (string)(value), err.Msg)
		}
	}

	// decode the json value for each message received into an Order Structure
	jsonDecode := func(messages <-chan []byte, orders chan<- trading_engine.Order) {
		for msg := range messages {
			var order trading_engine.Order
			order.FromJSON(msg)
			orders <- order
		}
	}

	// process each order by the trading engine and forward trades to the trades channel
	processOrders := func(engine *trading_engine.TradingEngine, orders <-chan trading_engine.Order, tradeChan chan<- []trading_engine.Trade, n int) {
		for order := range orders {
			trades := engine.Process(order)
			ordersCompleted++
			tradesCompleted += len(trades)
			if len(trades) > 0 {
				tradeChan <- trades
			}
			if ordersCompleted >= n {
				close(tradeChan)
				done <- true
				return
			}
		}
	}

	publishTrades := func(tradeChan <-chan []trading_engine.Trade, finishTrades chan bool, closeChan bool) {
		// buffer the writes to the publisher by 10000 records at a time
		for trades := range tradeChan {
			for _, trade := range trades {
				rawTrade, _ := trade.ToJSON() // @todo thread error on encoding json object (low priority)
				producer.Input() <- rawTrade
			}
		}
		if closeChan {
			finishTrades <- true
		}
	}

	startTime := time.Now().UnixNano()
	benchmark.ResetTimer()

	go receiveProducerErrors()
	go receiveMessages(messages, benchmark.N)
	go jsonDecode(messages, orders)
	go processOrders(engine, orders, tradeChan, benchmark.N)
	go publishTrades(tradeChan, finishTrades, true)

	<-done
	<-finishTrades
	producer.Close()
	endTime := time.Now().UnixNano()
	timeout := (float64)(float64(time.Nanosecond) * float64(endTime-startTime) / float64(time.Second))
	fmt.Printf(
		"Total Orders: %d\n"+
			"Total Trades Generated: %d\n"+
			"Orders/second: %f\n"+
			"Trades Generated/second: %f\n"+
			"Pending Buy: %d\n"+
			"Lowest Ask: %f\n"+
			"Pending Sell: %d\n"+
			"Highest Bid: %f\n"+
			"Duration (seconds): %f\n\n",
		ordersCompleted,
		tradesCompleted,
		float64(ordersCompleted)/timeout,
		float64(tradesCompleted)/timeout,
		engine.OrderBook.PricePoints.Len(),
		engine.OrderBook.LowestAsk,
		engine.OrderBook.PricePoints.Len(),
		engine.OrderBook.HighestBid,
		timeout,
	)
}
