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
			id := "" //fmt.Sprintf("%d", rand.Int())
			rnd := rand.Float64()
			price := 4000100.00 - 1000000*rnd
			amount := 10000.0 - 9000*rand.Float64()
			var side trading_engine.OrderSide
			if rand.Intn(2) == 1 {
				side = trading_engine.BUY
			} else {
				side = trading_engine.SELL
			}
			order := &trading_engine.Order{
				ID:       id,
				Side:     side,
				Category: 1,
				Amount:   amount,
				Price:    price,
			}
			tradingEngine.Process(order)
		}
	})
	endTime := time.Now().UnixNano()
	timeout := (float64)(int64(time.Nanosecond) * (endTime - startTime) / int64(time.Second))
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
