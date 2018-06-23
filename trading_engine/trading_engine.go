package trading_engine

// TradingEngine contains the current order book and information about the service since it was created
type TradingEngine interface {
	Process(order Order) []Trade
	GetOrderBook() OrderBook
}

type tradingEngine struct {
	OrderBook OrderBook
	Symbol    string
}

// NewTradingEngine creates a new trading engine that contains an empty order book and can start receving requests
func NewTradingEngine() TradingEngine {
	orderBook := NewOrderBook()
	return &tradingEngine{
		OrderBook: orderBook,
	}
}

// Process a single order and returned all the trades that can be satisfied instantly
func (engine *tradingEngine) Process(order Order) []Trade {
	return engine.OrderBook.Process(order)
}

func (engine tradingEngine) GetOrderBook() OrderBook {
	return engine.OrderBook
}
