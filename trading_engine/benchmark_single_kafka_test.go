package trading_engine_test

import (
	"fmt"
	"log"
	"testing"
	"time"
	"trading_engine/net"
	"trading_engine/trading_engine"
)

func BenchmarkSingleKafkaTest(benchmark *testing.B) {
	kafkaBroker := "kafka:9092"
	kafkaOrderTopic := "trading.order.btc.eth"
	kafkaOrderConsumer := "benchmark_single_kafka_test_1"
	consumer := net.NewKafkaConsumer([]string{kafkaBroker}, []string{kafkaOrderTopic})
	consumer.Start(kafkaOrderConsumer)
	defer consumer.Close()

	// generateOrdersInKafka(benchmark.N)
	engine := trading_engine.NewTradingEngine()

	orderChan := make(chan trading_engine.Order, 1000)
	defer close(orderChan)
	stopChan := make(chan int, 100)
	defer close(stopChan)
	doneChan := make(chan bool, 5)
	defer close(doneChan)
	processChan := make(chan bool, 5)
	defer close(processChan)
	startTime := time.Now().UnixNano()
	benchmark.ResetTimer()

	go listenForOrders(consumer, orderChan, benchmark.N, "1", stopChan)
	// go listenForOrders(consumer, orderChan, benchmark.N, "1", stopChan)
	// go listenForOrders(consumer, orderChan, benchmark.N, "1", stopChan)
	// go listenForOrders(consumer, orderChan, benchmark.N, "1", stopChan)

	go processOrders(engine, doneChan, processChan, orderChan)

	nr := <-stopChan
	log.Println(nr)
	processChan <- true
	<-doneChan

	endTime := time.Now().UnixNano()
	timeout := (float64(time.Nanosecond) * float64(endTime-startTime) / float64(time.Second))
	duration := (float64(time.Nanosecond) * float64(endTime-startTime) / float64(time.Millisecond))
	fmt.Printf(
		"Total Orders: %d\n"+
			"Total Trades: %d\n"+
			"Orders/second: %f\n"+
			"Trades/second: %f\n"+
			"Pending Buy: %d\n"+
			"Lowest Ask: %f\n"+
			"Pending Sell: %d\n"+
			"Highest Bid: %f\n"+
			"Duration (ms): %f\n\n",
		engine.OrdersCompleted,
		engine.TradesCompleted,
		float64(engine.OrdersCompleted)/timeout,
		float64(engine.TradesCompleted)/timeout,
		engine.OrderBook.PricePoints.Len(),
		engine.OrderBook.LowestAsk,
		engine.OrderBook.PricePoints.Len(),
		engine.OrderBook.HighestBid,
		duration,
	)
}

func listenForOrders(consumer *net.KafkaConsumer, orderChan chan trading_engine.Order, n int, index string, stopChan chan int) {
	messageChan := consumer.GetMessageChan()
	i := 0
	for msg := range messageChan {
		consumer.MarkOffset(msg, "")
		var order trading_engine.Order
		fmt.Sscanf((string)(msg.Value), "order[%d][ %s ] %d - %f @ %f", &order.Category, &order.ID, &order.Side, &order.Amount, &order.Price)
		orderChan <- order
		i++
		if i >= n {
			break
		}
	}
	stopChan <- i
}

// func process(engine *trading_engine.TradingEngine, order *trading_engine.Order, w *sync.WaitGroup) {
// 	engine.Process(order)
// 	w.Done()
// }

func processOrders(engine *trading_engine.TradingEngine, doneChan chan bool, processChan chan bool, orderChan chan trading_engine.Order) {
	for {
		select {
		case order := <-orderChan:
			{
				engine.Process(order)
			}
		case <-processChan:
			{
				doneChan <- true
				return
			}
		}
		// for i := range trades {
		// 	trade := trades[i]
		// 	tradeMsg := net.EncodeTrade(trade.MakerOrderID, trade.TakerOrderID, trade.Price, trade.Amount)
		// 	producer.Send(tradeMsg)
		// }
	}
}

// func generateOrders(n int) {
// 	rand.Seed(42)
// 	kafkaBroker := "kafka:9092"
// 	kafkaOrderTopic := "trading.order.btc.eth"

// 	producer := net.NewKafkaProducer([]string{kafkaBroker}, kafkaOrderTopic)
// 	err := producer.Start()
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	for i := 0; i < n; i++ {
// 		rnd := rand.Float64()
// 		msg := fmt.Sprintf("order[%d][ %s ] %d - %f @ %f", 1, "id", 1+rand.Intn(2)%2, 10001.0-10000*rand.Float64(), 4000100.00-(float64)(i)-1000000*rnd)
// 		_, _, err := producer.Send(([]byte)(msg))
// 		if err != nil {
// 			log.Println(err)
// 			break
// 		}
// 	}
// 	producer.Close()
// }
