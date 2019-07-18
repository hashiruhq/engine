package benchmarks_test

// import (
// 	"fmt"
// 	"testing"
// 	"time"

// 	"gitlab.com/around25/products/matching-engine/engine"
// 	"gitlab.com/around25/products/matching-engine/model"
// 	"gitlab.com/around25/products/matching-engine/net"
// )

// func BenchmarkKafkaConsumer(benchmark *testing.B) {
// 	ngin := engine.NewTradingEngine("btcusd", 8, 8)

// 	// generateOrdersInKafka(benchmark.N)

// 	kafkaBroker := "kafka:9092"
// 	kafkaOrderTopic := "trading.order.btc.eth"
// 	kafkaOrderConsumer := "benchmark_kafka_consumer_test"
// 	ordersCompleted := 0
// 	tradesCompleted := 0

// 	consumer := net.NewKafkaPartitionConsumer(kafkaOrderConsumer, []string{kafkaBroker}, []string{kafkaOrderTopic})
// 	consumer.Start()
// 	defer consumer.Close()

// 	orders := make(chan model.Order, 10000)
// 	defer close(orders)

// 	messages := make(chan []byte, 10000)
// 	defer close(messages)

// 	done := make(chan bool)
// 	defer close(done)

// 	decode := func(messages <-chan []byte, orders chan<- model.Order) {
// 		for {
// 			msg, more := <-messages
// 			if !more {
// 				return
// 			}
// 			var order model.Order
// 			order.FromBinary(msg)
// 			orders <- order
// 		}
// 	}

// 	receiveMessages := func(messages chan<- []byte, n int) {
// 		msgChan := consumer.GetMessageChan()
// 		for j := 0; j < n; j++ {
// 			msg := <-msgChan
// 			// consumer.MarkOffset(msg, "")
// 			messages <- msg.Value
// 		}
// 	}

// 	startTime := time.Now().UnixNano()
// 	benchmark.ResetTimer()

// 	go receiveMessages(messages, benchmark.N)
// 	go decode(messages, orders)
// 	go func(ngin engine.TradingEngine, orders <-chan model.Order, n int) {
// 		events := make([]model.Event, 1000)
// 		for {
// 			order := <-orders
// 			ngin.Process(order, &events)
// 			ordersCompleted++
// 			if len(events) >= cap(events)/2+1 {
// 				eventsCompleted += len(events)
// 				events = events[0:0]
// 			}
// 			if ordersCompleted >= n {
// 				done <- true
// 				return
// 			}
// 		}
// 	}(ngin, orders, benchmark.N)

// 	<-done
// 	endTime := time.Now().UnixNano()
// 	timeout := (float64)(float64(time.Nanosecond) * float64(endTime-startTime) / float64(time.Second))
// 	fmt.Printf(
// 		"Total Orders: %d\n"+
// 			"Total Eventss: %d\n"+
// 			"Orders/second: %f\n"+
// 			"Events/second: %f\n"+
// 			"Pending Buy: %d\n"+
// 			"Lowest Ask: %d\n"+
// 			"Pending Sell: %d\n"+
// 			"Highest Bid: %d\n"+
// 			"Duration (seconds): %f\n\n",
// 		ordersCompleted,
// 		eventsCompleted,
// 		float64(ordersCompleted)/timeout,
// 		float64(eventsCompleted)/timeout,
// 		ngin.GetOrderBook().GetMarket()[0].Len(),
// 		ngin.GetOrderBook().GetLowestAsk(),
// 		ngin.GetOrderBook().GetMarket()[1].Len(),
// 		ngin.GetOrderBook().GetHighestBid(),
// 		timeout,
// 	)
// }
