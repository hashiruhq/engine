package benchmarks_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"gitlab.com/around25/products/matching-engine/engine"
// 	"gitlab.com/around25/products/matching-engine/model"
// 	"gitlab.com/around25/products/matching-engine/net"
// )

// func init() {
// 	GenerateOrdersInKafka(KAFKA_BROKER, KAFKA_CONSUMER_TOPIC, 2000000)
// }

// func BenchmarkKafkaConsumer(benchmark *testing.B) {
// 	ordersCompleted := 0
// 	eventsCompleted := 0

// 	ngin := engine.NewTradingEngine(KAFKA_CONSUMER_MARKET, 8, 8)
// 	readerConfig := net.KafkaReaderConfig{}

// 	consumer := net.NewKafkaConsumer(readerConfig, []string{KAFKA_BROKER}, KAFKA_CONSUMER_TOPIC, 0)
// 	consumer.SetOffset(-2) // oldest offset available
// 	consumer.Start(context.Background())
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
// 			msg, more := <-msgChan
// 			if !more {
// 				return
// 			}
// 			messages <- msg.Value
// 		}
// 	}

// 	startTime := time.Now().UnixNano()
// 	benchmark.ResetTimer()

// 	go receiveMessages(messages, benchmark.N)
// 	go decode(messages, orders)
// 	go func(ngin engine.TradingEngine, orders <-chan model.Order, n int) {
// 		events := make([]model.Event, 0, 1000)
// 		for {
// 			order, more := <-orders
// 			if !more {
// 				return
// 			}
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

// 	select {
// 	case <-done:
// 	case <-time.After(3 * time.Second):
// 		fmt.Println("FAIL: Test timeed out after 3 seconds")
// 	}

// 	endTime := time.Now().UnixNano()
// 	timeout := (float64)(float64(time.Nanosecond) * float64(endTime-startTime) / float64(time.Second))
// 	fmt.Printf(
// 		"Total Orders: %d\n"+
// 			"Total Events: %d\n"+
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
