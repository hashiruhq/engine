package main

import (
	"runtime"
	"sync"
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
	tradingEngine := trading_engine.NewTradingEngine()

	// w.Add(runtime.NumCPU())
	// for i := 0; i < runtime.NumCPU(); i++ {
	// 	go processTrades(i, tradePool)
	// }

	w.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go generateOrders(tradingEngine)
	}

	w.Wait()
}

func generateOrders(tradingEngine *trading_engine.TradingEngine) {
	// order := &trading_engine.Order{
	// 	ID:       id,
	// 	Side:     side,
	// 	Category: 1,
	// 	Amount:   amount,
	// 	Price:    price,
	// }
	// tradingEngine.Process(order)
	w.Done()
}
