package engine

// TradingEngine contains the current order book and information about the service since it was created
type TradingEngine interface {
	Process(order Order) []Trade
	GetOrderBook() OrderBook
	BackupMarket() MarketBackup
	LoadMarket(MarketBackup) error
	CancelOrder(order Order) bool
	ProcessEvent(order Order) interface{}
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
func (ngin *tradingEngine) Process(order Order) []Trade {
	return ngin.OrderBook.Process(order)
}

func (ngin *tradingEngine) CancelOrder(order Order) bool {
	return ngin.OrderBook.Cancel(order)
}

func (ngin *tradingEngine) LoadMarket(market MarketBackup) error {
	return ngin.GetOrderBook().Load(market)
}

func (ngin *tradingEngine) BackupMarket() MarketBackup {
	return ngin.GetOrderBook().Backup()
}

func (ngin *tradingEngine) ProcessEvent(order Order) interface{} {
	switch order.EventType {
	case CommandType_NewOrder:
		return ngin.Process(order)
	case CommandType_CancelOrder:
		return ngin.CancelOrder(order)
	case CommandType_BackupMarket:
		return ngin.BackupMarket()
	}
	return nil
}

func (ngin tradingEngine) GetOrderBook() OrderBook {
	return ngin.OrderBook
}
