package trading_engine

import "sync"

// TradingEngine contains the current order book and information about the service since it was created
type TradingEngine struct {
	OrderBook       OrderBook
	OrdersCompleted int64
	TradesCompleted int64
	lock            sync.Mutex
	Symbol          string
	LastRestart     string
}

// NewTradingEngine creates a new trading engine that contains an empty order book and can start receving requests
func NewTradingEngine() *TradingEngine {
	orderBook := NewOrderBook()
	return &TradingEngine{OrderBook: orderBook, OrdersCompleted: 0, TradesCompleted: 0}
}

// Process a single order and returned all the trades that can be satisfied instantly
func (engine *TradingEngine) Process(order Order) []Trade {
	engine.lock.Lock()
	trades := engine.OrderBook.Process(order)
	engine.OrdersCompleted++
	engine.TradesCompleted += (int64)(len(trades))
	engine.lock.Unlock()
	return trades
}
