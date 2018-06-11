package trading_engine

import (
	"errors"
	"fmt"
)

// Start the engine by consuming orders from the order pool and send them to the tradePool
func (engine *TradingEngine) Start(orderPool <-chan *Order, tradePool chan<- *Trade) error {
	for order := range orderPool {
		engine.OrdersCompleted++
		trades, err := engine.OrderBook.Process(order)
		if err == nil {
			engine.TradesCompleted += (int64)(len(trades))
			for index := range trades {
				select {
				case tradePool <- trades[index]:
				default:
					fmt.Println("Unable to send results", len(tradePool), cap(tradePool))
				}
			}
		} else {
			fmt.Println("Error executing order", err, order)
		}
	}
	return errors.New("ERROR: Stopping engine due to closing order channel")
}

func (engine *TradingEngine) Process(order *Order) []*Trade {
	trades, err := engine.OrderBook.Process(order)
	engine.mtx.Lock()
	engine.OrdersCompleted++
	if err == nil {
		engine.TradesCompleted += (int64)(len(trades))
	} else {
		fmt.Println("Error executing order", err, order)
	}
	engine.mtx.Unlock()
	return trades
}
