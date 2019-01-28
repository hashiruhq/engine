package engine

// TradingEngine contains the current order book and information about the service since it was created
type TradingEngine interface {
	Process(order Order) []Trade
	GetOrderBook() OrderBook
	BackupMarket() Market
	LoadMarket(Market) error
	CancelOrder(order Order) *BookEntry
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

func (ngin *tradingEngine) CancelOrder(order Order) *BookEntry {
	return ngin.OrderBook.Cancel(order.ID)
}

func (ngin *tradingEngine) LoadMarket(market Market) error {
	return ngin.GetOrderBook().Load(market)
}

func (ngin *tradingEngine) BackupMarket() Market {
	return ngin.GetOrderBook().Backup()
}

func (ngin *tradingEngine) ProcessEvent(order Order) interface{} {
	switch order.EventType {
	case EventTypeNewOrder:
		return ngin.Process(order)
	case EventTypeCancelOrder:
		return ngin.CancelOrder(order)
	case EventTypeBackupMarket:
		return ngin.BackupMarket()
	}
	return nil
}

func (ngin tradingEngine) GetOrderBook() OrderBook {
	return ngin.OrderBook
}
