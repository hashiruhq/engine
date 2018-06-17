package trading_engine_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"trading_engine/trading_engine"
)

func BenchmarkWithRandomData(benchmark *testing.B) {
	rand.Seed(42)
	tradingEngine := trading_engine.NewTradingEngine()
	startTime := time.Now().UnixNano()
	benchmark.ResetTimer()
	benchmark.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rnd := rand.Float64()
			order := trading_engine.Order{
				Price:    4000100.00 - 1000000*rnd,
				Amount:   10001.0 - 10000*rand.Float64(),
				Side:     rand.Intn(2) % 2,
				ID:       "",
				Category: trading_engine.LIMIT_ORDER,
			}
			tradingEngine.Process(order)
		}
	})
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
		tradingEngine.OrdersCompleted,
		tradingEngine.TradesCompleted,
		float64(tradingEngine.OrdersCompleted)/timeout,
		float64(tradingEngine.TradesCompleted)/timeout,
		tradingEngine.OrderBook.PricePoints.Len(),
		tradingEngine.OrderBook.LowestAsk,
		tradingEngine.OrderBook.PricePoints.Len(),
		tradingEngine.OrderBook.HighestBid,
		timeout,
	)
}
