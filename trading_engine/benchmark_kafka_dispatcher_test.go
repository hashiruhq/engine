package trading_engine_test

import (
	"fmt"
	"testing"
	"time"
	"trading_engine/config"
	"trading_engine/net"
	"trading_engine/queue"
	"trading_engine/trading_engine"
)

// BenchmarkKafkaDispatcher test speed using dispatcher from
// this article http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/

func BenchmarkKafkaDispatcher(benchmark *testing.B) {
	kafkaBroker := "kafka:9092"
	kafkaOrderTopic := "trading.order.btc.eth"
	kafkaOrderConsumer := "benchmark_kafka_dispatcher_test"
	consumer := net.NewKafkaConsumer([]string{kafkaBroker}, []string{kafkaOrderTopic})
	consumer.Start(kafkaOrderConsumer)
	defer consumer.Close()
	const MAX_RECORDS = 240000
	// MAX_RECORDS := benchmark.N

	// generateOrders(MAX_RECORDS)
	engine := trading_engine.NewTradingEngine()

	dispatcher := queue.NewDispatcher(engine, config.Config.GetInt("MAX_WORKERS"))
	dispatcher.Run()

	startTime := time.Now().UnixNano()
	benchmark.ResetTimer()

	listenForOrders := func(engine *trading_engine.TradingEngine, consumer *net.KafkaConsumer, n int) {
		messageChan := consumer.GetMessageChan()
		i := 0
		for msg := range messageChan {
			consumer.MarkOffset(msg, "")
			order := trading_engine.Order{}
			fmt.Sscanf((string)(msg.Value), "order[%d][ %s ] %d - %f @ %f", &order.Category, &order.ID, &order.Side, &order.Amount, &order.Price)
			job := queue.Job{Order: order}
			queue.JobQueue <- job
			// engine.Process(&order)
			i++
			if i >= n {
				break
			}
		}
	}

	go listenForOrders(engine, consumer, MAX_RECORDS/6)
	go listenForOrders(engine, consumer, MAX_RECORDS/6)
	go listenForOrders(engine, consumer, MAX_RECORDS/6)
	go listenForOrders(engine, consumer, MAX_RECORDS/6)
	go listenForOrders(engine, consumer, MAX_RECORDS/6)
	go listenForOrders(engine, consumer, MAX_RECORDS/6)

	time.Sleep(2 * time.Second)
	dispatcher.Stop()

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
