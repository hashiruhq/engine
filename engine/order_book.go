package engine

import (
	"gitlab.com/around25/products/matching-engine/model"
)

// OrderBook interface
// Defines what constitudes an order book and how we can interact with it
type OrderBook interface {
	Process(model.Order, *[]model.Event)
	Cancel(model.Order, *[]model.Event)
	GetLastTradePriceFromEvents(*[]model.Event) uint64
	ActivateStopOrders(price uint64, events *[]model.Event) *[]model.Order
	GetHighestBid() uint64
	GetLowestAsk() uint64
	Load(model.MarketBackup) error
	Backup() model.MarketBackup
	GetMarket() []*SkipList
	GetMarketID() string
	GetPricePrecision() int
	GetVolumePrecision() int
}

type orderBook struct {
	MarketID        string
	PricePrecision  int
	VolumePrecision int

	// unfilled limit orders
	BuyEntries  *SkipList
	SellEntries *SkipList
	LowestAsk   uint64
	HighestBid  uint64

	// market orders waiting liquidity
	BuyMarketEntries  []model.Order
	SellMarketEntries []model.Order

	// pending stop loss orders
	StopEntryOrders  *SkipList
	StopLossOrders   *SkipList
	LowestEntryPrice uint64
	HighestLossPrice uint64
}

// NewOrderBook Creates a new empty order book for the trading engine
func NewOrderBook(marketID string, pricePrecision, volumePrecision int) OrderBook {
	return &orderBook{
		MarketID:        marketID,
		PricePrecision:  pricePrecision,
		VolumePrecision: volumePrecision,
		// Limit order data
		LowestAsk:   0,
		HighestBid:  0,
		BuyEntries:  NewPricePoints(),
		SellEntries: NewPricePoints(),
		// Market order data
		BuyMarketEntries:  make([]model.Order, 0, 0),
		SellMarketEntries: make([]model.Order, 0, 0),
		// Stop order data
		StopEntryOrders:  NewPricePoints(),
		StopLossOrders:   NewPricePoints(),
		LowestEntryPrice: 0,
		HighestLossPrice: 0,
	}
}

func (book orderBook) GetMarketID() string {
	return book.MarketID
}

func (book orderBook) GetPricePrecision() int {
	return book.PricePrecision
}

func (book orderBook) GetVolumePrecision() int {
	return book.VolumePrecision
}

// GetHighestBid returns the highest bid of the current market
func (book orderBook) GetHighestBid() uint64 {
	return book.HighestBid
}

// GetLowestAsk returns the lowest ask of the current market
func (book orderBook) GetLowestAsk() uint64 {
	return book.LowestAsk
}

// GetHighestLossPrice returns the highest price for which there is a stop loss order
func (book orderBook) GetHighestLossPrice() uint64 {
	return book.HighestLossPrice
}

// GetLowestEntryPrice return the minimum price for which there is a stop entry order
func (book orderBook) GetLowestEntryPrice() uint64 {
	return book.LowestEntryPrice
}

// GetMarket returns the list of price points for the market
func (book orderBook) GetMarket() []*SkipList {
	return []*SkipList{book.BuyEntries, book.SellEntries}
}

// Process a new received order and return a list of events make
func (book *orderBook) Process(order model.Order, events *[]model.Event) {
	switch order.EventType {
	case model.CommandType_NewOrder:
		book.processOrder(order, events)
		// process activated stop orders
		if events != nil && len(*events) > 0 {
			lastPrice := book.GetLastTradePriceFromEvents(events)
			stopOrders := book.ActivateStopOrders(lastPrice, events)
			if stopOrders != nil && len(*stopOrders) > 0 {
				for _, stopOrder := range *stopOrders {
					stopOrder.Stop = model.StopLoss_None
					// recursively process the activated order as a normal order
					book.Process(stopOrder, events)
				}
			}
		}
	case model.CommandType_CancelOrder:
		book.Cancel(order, events)
	}
}

// Process a new order and return a list of events resulted from the exchange
func (book *orderBook) processOrder(order model.Order, events *[]model.Event) {
	// if the order is a stop order then add it to the list of pending stop orders
	if order.Stop != model.StopLoss_None && order.StopPrice != 0 {
		book.addStopOrder(order)
		return
	}

	// for limit orders first process the limit order with the orderbook since you
	// can't have a pending market order and not have an empty order book
	if order.Type == model.OrderType_Limit && order.Side == model.MarketSide_Buy {
		book.processLimitBuy(order, events)
		// if there are sell market orders pending then keep loading the first one until there are
		// no more pending market orders or the limit order was filled.
		sellMarketCount := len(book.SellMarketEntries)
		if sellMarketCount == 0 {
			return
		}
		for i := 0; i < sellMarketCount; i++ {
			mko := book.popMarketSellOrder()
			mkOrder := book.processMarketSell(*mko, events)
			if mkOrder.Status != model.OrderStatus_Filled {
				book.lpushMarketSellOrder(mkOrder)
				break
			}
		}
		return
	}
	// similar to the above for sell limit orders
	if order.Type == model.OrderType_Limit && order.Side == model.MarketSide_Sell {
		book.processLimitSell(order, events)
		buyMarketCount := len(book.BuyMarketEntries)
		if buyMarketCount == 0 {
			return
		}
		for i := 0; i < buyMarketCount; i++ {
			mko := book.popMarketBuyOrder()
			mkOrder := book.processMarketBuy(*mko, events)
			if mkOrder.Status != model.OrderStatus_Filled {
				book.lpushMarketBuyOrder(mkOrder)
				break
			}
		}
		return
	}

	// market order either get filly filled or they get added to the pending market order list
	if order.Type == model.OrderType_Market && order.Side == model.MarketSide_Buy {
		if len(book.BuyMarketEntries) > 0 {
			book.pushMarketBuyOrder(order)
		} else {
			order := book.processMarketBuy(order, events)
			if order.Status != model.OrderStatus_Filled {
				book.pushMarketBuyOrder(order)
			}
			return
		}
	}
	// exactly the same for sell market orders
	if order.Type == model.OrderType_Market && order.Side == model.MarketSide_Sell {
		if len(book.SellMarketEntries) > 0 {
			book.pushMarketSellOrder(order)
		} else {
			order := book.processMarketSell(order, events)
			if order.Status != model.OrderStatus_Filled {
				book.pushMarketSellOrder(order)
			}
			return
		}
	}
}

// Cancel an order from the order book based on the order price and ID
func (book *orderBook) Cancel(order model.Order, events *[]model.Event) {
	// cancel stop orders
	if order.Stop != model.StopLoss_None {
		book.cancelStopOrder(order, events)
		return
	}
	// cancel non stop orders
	switch order.Type {
	case model.OrderType_Limit:
		book.cancelLimitOrder(order, events)
	case model.OrderType_Market:
		book.cancelMarketOrder(order, events)
	default:
	}
}

// Cancel a stop order based on a given order ID and set price
func (book *orderBook) cancelStopOrder(order model.Order, events *[]model.Event) {
	if order.Stop == model.StopLoss_None {
		return
	}
	price := order.StopPrice

	if order.Stop == model.StopLoss_Loss {
		if pricePoint, ok := book.StopLossOrders.Get(price); ok {
			for i := 0; i < len(pricePoint.Entries); i++ {
				if pricePoint.Entries[i].ID == order.ID {
					ord := pricePoint.Entries[i]
					ord.SetStatus(model.OrderStatus_Cancelled)
					*events = append(*events, model.NewOrderStatusEvent(book.MarketID, ord.ID, ord.Amount, ord.Funds, ord.Status))
					book.StopLossOrders.removeEntryByPriceAndIndex(price, pricePoint, i)
					return
				}
			}
		}
		return
	}
	if pricePoint, ok := book.StopEntryOrders.Get(order.Price); ok {
		for i := 0; i < len(pricePoint.Entries); i++ {
			if pricePoint.Entries[i].ID == order.ID {
				ord := pricePoint.Entries[i]
				ord.SetStatus(model.OrderStatus_Cancelled)
				*events = append(*events, model.NewOrderStatusEvent(book.MarketID, ord.ID, ord.Amount, ord.Funds, ord.Status))
				book.StopEntryOrders.removeEntryByPriceAndIndex(price, pricePoint, i)
				return
			}
		}
	}
	// if nothing was cancelled is possible the order was already activated and we should cancel it from the market
	order.Stop = model.StopLoss_None
	book.Cancel(order, events)
}

// Cancel a limit order based on a given order ID and set price
func (book *orderBook) cancelLimitOrder(order model.Order, events *[]model.Event) {
	if order.Side == model.MarketSide_Buy {
		if pricePoint, ok := book.BuyEntries.Get(order.Price); ok {
			for i := 0; i < len(pricePoint.Entries); i++ {
				if pricePoint.Entries[i].ID == order.ID {
					ord := pricePoint.Entries[i]
					ord.SetStatus(model.OrderStatus_Cancelled)
					*events = append(*events, model.NewOrderStatusEvent(book.MarketID, ord.ID, ord.Amount, ord.Funds, ord.Status))
					book.removeBuyBookEntry(ord.Price, pricePoint, i)
					// adjust highest bid
					if len(pricePoint.Entries) == 0 && book.HighestBid == ord.Price {
						iterator := book.BuyEntries.Seek(book.HighestBid)
						book.closeBidIterator(iterator)
					}
					return
				}
			}
		}
		return
	}
	if pricePoint, ok := book.SellEntries.Get(order.Price); ok {
		for i := 0; i < len(pricePoint.Entries); i++ {
			if pricePoint.Entries[i].ID == order.ID {
				ord := pricePoint.Entries[i]
				ord.SetStatus(model.OrderStatus_Cancelled)
				*events = append(*events, model.NewOrderStatusEvent(book.MarketID, ord.ID, ord.Amount, ord.Funds, ord.Status))
				book.removeSellBookEntry(ord.Price, pricePoint, i)
				// adjust lowest ask
				if len(pricePoint.Entries) == 0 && book.LowestAsk == ord.Price {
					iterator := book.SellEntries.Seek(book.LowestAsk)
					book.closeAskIterator(iterator)
				}
				return
			}
		}
	}
}

// Cancel a market order based on a given order ID
func (book *orderBook) cancelMarketOrder(order model.Order, events *[]model.Event) {
	if order.Side == model.MarketSide_Buy {
		for i := 0; i < len(book.BuyMarketEntries); i++ {
			if order.ID == book.BuyMarketEntries[i].ID {
				ord := book.BuyMarketEntries[i]
				ord.SetStatus(model.OrderStatus_Cancelled)
				*events = append(*events, model.NewOrderStatusEvent(book.MarketID, ord.ID, ord.Amount, ord.Funds, ord.Status))
				book.BuyMarketEntries = append(book.BuyMarketEntries[:i], book.BuyMarketEntries[i+1:]...)
				return
			}
		}
		return
	}
	for i := 0; i < len(book.SellMarketEntries); i++ {
		if order.ID == book.SellMarketEntries[i].ID {
			ord := book.SellMarketEntries[i]
			ord.SetStatus(model.OrderStatus_Cancelled)
			*events = append(*events, model.NewOrderStatusEvent(book.MarketID, ord.ID, ord.Amount, ord.Funds, ord.Status))
			book.SellMarketEntries = append(book.SellMarketEntries[:i], book.SellMarketEntries[i+1:]...)
			return
		}
	}
}

// Add a new book entry in the order book
// If the price point already exists then the book entry is simply added at the end of the pricepoint
// If the price point does not exist yet it will be created
func (book *orderBook) addBuyBookEntry(order model.Order) {
	book.BuyEntries.addOrder(order.Price, order)
}

func (book *orderBook) addSellBookEntry(order model.Order) {
	book.SellEntries.addOrder(order.Price, order)
}

// Remove a book entry from the order book
// The method will also remove the price point entry if both book entry lists are empty
func (book *orderBook) removeBuyBookEntry(price uint64, pricePoint *PricePoint, index int) {
	book.BuyEntries.removeEntryByPriceAndIndex(price, pricePoint, index)
}

func (book *orderBook) removeSellBookEntry(price uint64, pricePoint *PricePoint, index int) {
	book.SellEntries.removeEntryByPriceAndIndex(price, pricePoint, index)
}

func (book *orderBook) pushMarketBuyOrder(order model.Order) {
	book.BuyMarketEntries = append(book.BuyMarketEntries, order)
}

func (book *orderBook) pushMarketSellOrder(order model.Order) {
	book.SellMarketEntries = append(book.SellMarketEntries, order)
}

func (book *orderBook) lpushMarketBuyOrder(order model.Order) {
	book.BuyMarketEntries = append([]model.Order{order}, book.BuyMarketEntries...)
}

func (book *orderBook) lpushMarketSellOrder(order model.Order) {
	book.SellMarketEntries = append([]model.Order{order}, book.SellMarketEntries...)
}

func (book *orderBook) popMarketSellOrder() *model.Order {
	if len(book.SellMarketEntries) == 0 {
		return nil
	}
	index := 0
	order := book.SellMarketEntries[0]
	book.SellMarketEntries = append(book.SellMarketEntries[:index], book.SellMarketEntries[index+1:]...)
	return &order
}

func (book *orderBook) popMarketBuyOrder() *model.Order {
	if len(book.BuyMarketEntries) == 0 {
		return nil
	}
	index := 0
	order := book.BuyMarketEntries[0]
	book.BuyMarketEntries = append(book.BuyMarketEntries[:index], book.BuyMarketEntries[index+1:]...)
	return &order
}

func (book *orderBook) closeBidIterator(iterator Iterator) {
	if iterator == nil {
		return
	}
	pricePoint := iterator.Value()
	if len(pricePoint.Entries) != 0 {
		iterator.Close()
		return
	}
	if ok := iterator.Previous(); ok {
		book.HighestBid = iterator.Key()
		iterator.Close()
		return
	}
	book.HighestBid = 0
	iterator.Close()
	return
}

func (book *orderBook) closeAskIterator(iterator Iterator) {
	if iterator == nil {
		return
	}
	pricePoint := iterator.Value()
	if len(pricePoint.Entries) != 0 {
		iterator.Close()
		return
	}
	if ok := iterator.Next(); ok {
		book.LowestAsk = iterator.Key()
		iterator.Close()
		return
	}
	book.LowestAsk = 0
	iterator.Close()
	return
}
