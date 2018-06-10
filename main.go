package main

import (
	"fmt"
	"math/rand"
	"runtime"
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

// var w sync.WaitGroup
var count int

func main() {
	rand.Seed(42)
	count = 0
	orderPool := make(chan *trading_engine.Order, 1000000)
	tradePool := make(chan *trading_engine.Trade, 3000000)

	tradingEngine := trading_engine.NewTradingEngine()
	// w.Add(1)
	// go func(orderPool <-chan *trading_engine.Order, tradePool chan<- *trading_engine.Trade) {
	// 	tradingEngine.Start(orderPool, tradePool)
	// }(orderPool, tradePool)
	startTime := time.Now().UnixNano() // / (int64(time.Millisecond)/int64(time.Nanosecond))

	// w.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go processTrades(i, tradePool)
	}
	// w.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU()/2; i++ {
		go generateOrders(orderPool)
	}

	go func(orderPool chan *trading_engine.Order, tradePool chan *trading_engine.Trade) {
		time.Sleep(10 * time.Second)
		fmt.Println("Closing pools")
		close(orderPool)
		// close(tradePool)
	}(orderPool, tradePool)

	tradingEngine.Start(orderPool, tradePool)

	close(tradePool)
	endTime := time.Now().UnixNano()

	timeout := (float64)(int64(time.Nanosecond) * (endTime - startTime) / int64(time.Second))
	// time.Sleep(30 * time.Second)
	fmt.Printf(
		"Total Orders: %d\n"+
			"Total Trades: %d\n"+
			"Orders/second: %f\n"+
			"Trades/second: %f\n"+
			"Pending Buy: %d\n"+
			"Lowest Bid: %f\n"+
			"Pending Sell: %d\n"+
			"Highest Ask: %f\n"+
			"Duration (seconds): %f",
		tradingEngine.OrdersCompleted,
		tradingEngine.TradesCompleted,
		float64(tradingEngine.OrdersCompleted)/timeout,
		float64(tradingEngine.TradesCompleted)/timeout,
		len(tradingEngine.OrderBook.BuyOrders),
		tradingEngine.OrderBook.BuyOrders[0].Price,
		len(tradingEngine.OrderBook.SellOrders),
		tradingEngine.OrderBook.SellOrders[0].Price,
		timeout,
	)
}

func generateOrders(orderPool chan<- *trading_engine.Order) {
	for {
		id := fmt.Sprintf("%d", rand.Int())
		price := rand.Float64() * 100
		amount := rand.Float64() * 10000
		var side trading_engine.OrderSide
		if rand.Intn(2) == 1 {
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
		select {
		case orderPool <- order:
		default:
			fmt.Println("Order not sent", order, len(orderPool), cap(orderPool))
			return
		}
	}
}

func processTrades(index int, tradePool <-chan *trading_engine.Trade) {
	for trade := range tradePool {
		if trade.Amount >= 0 {
			count++
		}
	}
	fmt.Println("Trade pool closed. Stop processing trades.")
}
