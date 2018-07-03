package trading_engine

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
func (engine *tradingEngine) Process(order Order) []Trade {
	return engine.OrderBook.Process(order)
}

func (engine *tradingEngine) CancelOrder(order Order) *BookEntry {
	return engine.OrderBook.Cancel(order.ID)
}

func (engine *tradingEngine) LoadMarket(market Market) error {
	return engine.GetOrderBook().Load(market)
}

func (engine *tradingEngine) BackupMarket() Market {
	return engine.GetOrderBook().Backup()
}

func (engine *tradingEngine) ProcessEvent(order Order) interface{} {
	switch order.EventType {
	case EventTypeNewOrder:
		return engine.Process(order)
	case EventTypeCancelOrder:
		return engine.CancelOrder(order)
	case EventTypeBackupMarket:
		return engine.BackupMarket()
	}
	return nil
}

func (engine tradingEngine) GetOrderBook() OrderBook {
	return engine.OrderBook
}
