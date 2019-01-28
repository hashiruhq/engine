package benchmarks_test

import (
	"log"
	"testing"

	"gitlab.com/around25/products/matching-engine/engine"
	"gitlab.com/around25/products/matching-engine/net"
	"gitlab.com/around25/products/matching-engine/queue"

	"github.com/Shopify/sarama"
)

const (
	kafkaBroker        = "kafka:9092"
	kafkaOrderTopic    = "trading.order.btc.eth"
	kafkaTradeTopic    = "trading.trade.btc.eth"
	kafkaOrderConsumer = "benchmark_kafka_consumer_producer_test"
)

func init() {
	// generateOrdersInKafka(500000)
}

func BenchmarkNetworkProcessChannels(benchmark *testing.B) {
	ngin := engine.NewTradingEngine()
	ordersCompleted := 0
	// tradesCompleted := 0
	producer := net.NewKafkaAsyncProducer([]string{kafkaBroker})
	producer.Start()
	consumer := net.NewKafkaPartitionConsumer(kafkaOrderConsumer, []string{kafkaBroker}, []string{kafkaOrderTopic})
	consumer.Start()
	defer consumer.Close()
	messages := make(chan []byte, 20000)
	defer close(messages)
	orders := make(chan engine.Order, 20000)
	defer close(orders)
	tradeChan := make(chan []engine.Trade, 20000)
	done := make(chan bool)
	defer close(done)
	finishTrades := make(chan bool)
	defer close(finishTrades)

	receiveMessages := func(messages chan<- []byte, n int) {
		msgChan := consumer.GetMessageChan()
		for j := 0; j < n; j++ {
			msg := <-msgChan
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

	jsonDecode := func(messages <-chan []byte, orders chan<- engine.Order) {
		for msg := range messages {
			var order engine.Order
			order.FromJSON(msg)
			orders <- order
		}
	}

	processOrders := func(ngin engine.TradingEngine, orders <-chan engine.Order, tradeChan chan<- []engine.Trade, n int) {
		for order := range orders {
			trades := ngin.Process(order)
			ordersCompleted++
			// tradesCompleted += len(trades)
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

	publishTrades := func(tradeChan <-chan []engine.Trade, finishTrades chan bool, closeChan bool) {
		for trades := range tradeChan {
			for _, trade := range trades {
				rawTrade, _ := trade.ToJSON() // @todo thread error on encoding json object (low priority)
				producer.Input() <- &sarama.ProducerMessage{
					Topic: kafkaTradeTopic,
					Value: sarama.ByteEncoder(rawTrade),
				}
			}
		}
		if closeChan {
			finishTrades <- true
		}
	}

	// startTime := time.Now().UnixNano()
	benchmark.ResetTimer()

	go receiveProducerErrors()
	go receiveMessages(messages, benchmark.N)
	go jsonDecode(messages, orders)
	go processOrders(ngin, orders, tradeChan, benchmark.N)
	go publishTrades(tradeChan, finishTrades, true)

	<-done
	<-finishTrades
	// benchmark.StopTimer()
	producer.Close()
}

func BenchmarkNetworkProcessEventRing(b *testing.B) {
	ngin := engine.NewTradingEngine()
	// ordersCompleted := 0
	// tradesCompleted := 0
	producer := net.NewKafkaAsyncProducer([]string{kafkaBroker})
	producer.Start()
	consumer := net.NewKafkaPartitionConsumer(kafkaOrderConsumer, []string{kafkaBroker}, []string{kafkaOrderTopic})
	consumer.Start()
	defer consumer.Close()
	done := make(chan bool)
	defer close(done)
	finishTrades := make(chan bool)
	defer close(finishTrades)

	msgQueue := queue.NewBuffer(1 << 15)
	orderQueue := queue.NewBuffer(1 << 15)
	tradeQueue := queue.NewBuffer(1 << 15)

	if b.N >= 2000000 {
		log.Fatalln("Limit reached")
	}

	go func(q *queue.Buffer, n int) {
		msgChan := consumer.GetMessageChan()
		for j := 0; j < n; j++ {
			msg := <-msgChan
			// consumer.MarkOffset(msg, "")
			q.Write(engine.NewEvent(msg))
		}
	}(msgQueue, b.N)

	go func(n int) {
		for i := 0; i < n; i++ {
			event := msgQueue.Read()
			event.Decode()
			orderQueue.Write(event)
		}
	}(b.N)

	go func(oq, tq *queue.Buffer, ngin engine.TradingEngine, n int) {
		for i := 0; i < n; i++ {
			event := oq.Read()
			event.SetTrades(ngin.Process(event.Order))
			tq.Write(event)
		}
	}(orderQueue, tradeQueue, ngin, b.N)

	receiveProducerErrors := func() {
		errors := producer.Errors()
		for err := range errors {
			value, _ := err.Msg.Value.Encode()
			log.Print("Error received from trades producer", (string)(value), err.Msg)
		}
	}
	// startTime := time.Now().UnixNano()
	// b.ResetTimer()
	go receiveProducerErrors()

	// go func(tq *queue.Gringo, n int) {
	for i := 0; i < b.N; i++ {
		event := tradeQueue.Read()
		for _, trade := range event.Trades {
			rawTrade, _ := trade.ToJSON()
			producer.Input() <- &sarama.ProducerMessage{
				Topic: kafkaTradeTopic,
				Value: sarama.ByteEncoder(rawTrade),
			}
		}
	}
	// 	done <- true
	// }(tradeQueue, b.N)
	// <-done
	// b.StopTimer()
	producer.Close()
	// endTime := time.Now().UnixNano()
}
