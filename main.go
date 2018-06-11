package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
	"trading_engine/trading_engine"
)

// User test structure
// type User struct {
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Age       int
// 	Languages []string `json:"languages"`
// }

// func printUserData(jsonData string, age int) string {
// 	var output string
// 	res := &User{}
// 	json.Unmarshal([]byte(jsonData), &res)
// 	if res.Age > age {
// 		output = fmt.Sprintf("User %s %s, who's %d can code in the following languages: %s\n", res.FirstName, res.LastName, res.Age, strings.Join(res.Languages, ", "))
// 	} else {
// 		output = fmt.Sprintf("User %s %s must be over %d before we can print their details", res.FirstName, res.LastName, age)
// 	}
// 	return output
// }

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var w sync.WaitGroup

func main() {
	rand.Seed(42)

	tradingEngine := trading_engine.NewTradingEngine()
	startTime := time.Now().UnixNano()

	// w.Add(runtime.NumCPU())
	// for i := 0; i < runtime.NumCPU(); i++ {
	// 	go processTrades(i, tradePool)
	// }

	w.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go generateOrders(tradingEngine)
	}

	w.Wait()

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
			"Duration (seconds): %f",
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

func generateOrders(tradingEngine *trading_engine.TradingEngine) {
	for i := 0; i < 100000; i++ {
		id := "" //fmt.Sprintf("%d", rand.Int())
		rnd := rand.Float64()
		price := 4000100.00 - float64(i) - 1000000*rnd
		amount := 10000.0 - 9000*rand.Float64()
		var side trading_engine.OrderSide
		if i%2 == 1 {
			side = trading_engine.BUY
		} else {
			side = trading_engine.SELL
		}
		if amount == 0 {
			continue
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
	w.Done()
}
