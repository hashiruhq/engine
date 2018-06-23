package trading_engine_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
	"trading_engine/trading_engine"
)

func BenchmarkWithRandomData(benchmark *testing.B) {
	rand.Seed(42)
	engine := trading_engine.NewTradingEngine()

	testFile := "/Users/cosmin/Incubator/go/src/trading_engine/priv/data/market.txt"

	// GenerateRandomRecordsInFile(&testFile, 2*benchmark.N)

	fh, err := os.Open(testFile)
	if err != nil {
		panic(err.Error())
	}
	defer fh.Close()
	bf := bufio.NewReader(fh)
	decoder := json.NewDecoder(bf)
	orders := make(chan trading_engine.Order, 10000)
	defer close(orders)
	done := make(chan bool)
	defer close(done)

	startTime := time.Now().UnixNano()
	benchmark.ResetTimer()
	go func(orders chan<- trading_engine.Order) {
		arr := make([]trading_engine.Order, 0, benchmark.N)
		for j := 0; j < benchmark.N; j++ {
			if !decoder.More() {
				break
			}
			order := trading_engine.Order{}
			decoder.Decode(&order)
			arr = append(arr, order)
		}
		startTime = time.Now().UnixNano()
		benchmark.ResetTimer()
		for _, order := range arr {
			orders <- order
		}
	}(orders)

	ordersCompleted := 0
	tradesCompleted := 0
	go func(orders <-chan trading_engine.Order, n int) {
		for {
			order := <-orders
			trades := engine.Process(order)
			ordersCompleted++
			tradesCompleted += len(trades)
			if ordersCompleted >= n {
				done <- true
				return
			}
		}
	}(orders, benchmark.N)

	<-done
	endTime := time.Now().UnixNano()
	timeout := (float64)(float64(time.Nanosecond) * float64(endTime-startTime) / float64(time.Second))
	fmt.Printf(
		"Total Orders: %d\n"+
			"Total Trades: %d\n"+
			"Orders/second: %f\n"+
			"Trades/second: %f\n"+
			"Pending Buy: %d\n"+
			"Lowest Ask: %f\n"+
			"Pending Sell: %d\n"+
			"Highest Bid: %f\n"+
			"Duration (seconds): %f\n\n",
		ordersCompleted,
		tradesCompleted,
		float64(ordersCompleted)/timeout,
		float64(tradesCompleted)/timeout,
		engine.GetOrderBook().PricePoints.Len(),
		engine.GetOrderBook().LowestAsk,
		engine.GetOrderBook().PricePoints.Len(),
		engine.GetOrderBook().HighestBid,
		timeout,
	)
}
