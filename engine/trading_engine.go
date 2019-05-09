package engine

// TradingEngine contains the current order book and information about the service since it was created
type TradingEngine interface {
	Process(order Order, trades *[]Trade)
	GetOrderBook() OrderBook
	BackupMarket() MarketBackup
	GetMarketDepth() *MarketDepth
	LoadMarket(MarketBackup) error
	CancelOrder(order Order) bool
	ProcessEvent(order Order, trades *[]Trade) interface{}
}

type tradingEngine struct {
	OrderBook OrderBook
	Symbol    string
}

// NewTradingEngine creates a new trading engine that contains an empty order book and can start receving requests
func NewTradingEngine(marketID string, pricePrecision, volumePrecision int) TradingEngine {
	orderBook := NewOrderBook(marketID, pricePrecision, volumePrecision)
	return &tradingEngine{
		OrderBook: orderBook,
	}
}

// Process a single order and returned all the trades that can be satisfied instantly
func (ngin *tradingEngine) Process(order Order, trades *[]Trade) {
	ngin.OrderBook.Process(order, trades)
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

func (ngin *tradingEngine) GetMarketDepth() *MarketDepth {
	return ngin.OrderBook.GetMarketDepth()
}

func (ngin *tradingEngine) ProcessEvent(order Order, trades *[]Trade) interface{} {
	switch order.EventType {
	case CommandType_NewOrder:
		ngin.Process(order, trades)
	case CommandType_CancelOrder:
		return ngin.CancelOrder(order)
	default:
		return nil
	}
	return nil
}

func (ngin tradingEngine) GetOrderBook() OrderBook {
	return ngin.OrderBook
}
