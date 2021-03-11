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
	ActivateStopOrders(events *[]model.Event) *[]model.Order
	GetHighestBid() uint64
	GetLowestAsk() uint64
	GetLastEventSeqID() uint64
	GetLastTradeSeqID() uint64
	Load(model.MarketBackup) error
	Backup() model.MarketBackup
	GetMarket() []*SkipList
	GetMarketID() string
	GetPricePrecision() int
	GetVolumePrecision() int
	GetHighestLossPrice() uint64
	GetLowestEntryPrice() uint64
	GetMarketOrders() ([]model.Order, []model.Order)
	AppendErrorEvent(*[]model.Event, model.ErrorCode, model.Order)
}

type orderBook struct {
	MarketID        string
	PricePrecision  int
	VolumePrecision int

	// sequence ids
	LastEventSeqID uint64
	LastTradeSeqID uint64

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
		// Sequence IDs
		LastEventSeqID: 0,
		LastTradeSeqID: 0,
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

func (book orderBook) GetLastEventSeqID() uint64 {
	return book.LastEventSeqID
}

func (book orderBook) GetLastTradeSeqID() uint64 {
	return book.LastTradeSeqID
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

// GetMarketOrders return the list of market orders
func (book orderBook) GetMarketOrders() ([]model.Order, []model.Order) {
	return book.BuyMarketEntries, book.SellMarketEntries
}

// Process a new received order and return a list of events make
func (book *orderBook) Process(order model.Order, events *[]model.Event) {
	switch order.EventType {
	case model.CommandType_NewOrder:
		// add acknowledgement event with status pending for stop orders and status untouched for limit/market orders
		order = book.ackOrder(order, events)
		// process the order normally
		book.processOrder(order, events)
		// process activated stop orders
		if events != nil && len(*events) > 0 {
			stopOrders := book.ActivateStopOrders(events)
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

func (book *orderBook) ackOrder(order model.Order, events *[]model.Event) model.Order {
	if order.Stop == model.StopLoss_None {
		order.Status = model.OrderStatus_Untouched
	}
	book.LastEventSeqID++
	event := model.NewOrderStatusEvent(book.LastEventSeqID, order.Market, order.Type, order.Side, order.ID, order.OwnerID, order.Price, order.Amount, order.Funds, order.Status, order.FilledAmount, order.UsedFunds)
	*events = append(*events, event)
	return order
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
		return
	}
	// similar to the above for sell limit orders
	if order.Type == model.OrderType_Limit && order.Side == model.MarketSide_Sell {
		book.processLimitSell(order, events)
		return
	}

	// market order either get filly filled or they get added to the pending market order list
	if order.Type == model.OrderType_Market && order.Side == model.MarketSide_Buy {
		_ = book.processMarketBuy(order, events)
		return
	}
	// exactly the same for sell market orders
	if order.Type == model.OrderType_Market && order.Side == model.MarketSide_Sell {
		_ = book.processMarketSell(order, events)
		return
	}
}

// Cancel an order from the order book based on the order price and ID
func (book *orderBook) Cancel(order model.Order, events *[]model.Event) {
	// cancel stop orders
	if order.Stop != model.StopLoss_None {
		book.cancelStopOrder(order, events)
		return
	}
	// ignore market orders since they are executed immediatly or cancelled
	// try to cancel limit order if any
	if order.Type == model.OrderType_Limit {
		if book.cancelLimitOrder(order, events) {
			return
		}
		book.AppendErrorEvent(events, model.ErrorCode_CancelFailed, order)
	}
}

// Cancel a stop loss order and return true if the order was found
func (book *orderBook) cancelStopLossOrder(order model.Order, events *[]model.Event) bool {
	price := order.StopPrice
	iterator := book.StopLossOrders.Seek(price)
	// price is outside the bounds of the list
	if iterator == nil {
		return false
	}
	// price is in the range but does not exist in the list
	if iterator.Key() != price {
		iterator.Close()
		return false
	}
	pricePoint := iterator.Value()
	for i := 0; i < len(pricePoint.Entries); i++ {
		if pricePoint.Entries[i].ID == order.ID {
			ord := pricePoint.Entries[i]
			ord.SetStatus(model.OrderStatus_Cancelled)
			book.LastEventSeqID++
			*events = append(*events, model.NewOrderStatusEvent(book.LastEventSeqID, book.MarketID, ord.Type, ord.Side, ord.ID, ord.OwnerID, ord.Price, ord.Amount, ord.Funds, ord.Status, ord.FilledAmount, ord.UsedFunds))
			book.StopLossOrders.removeEntryByPriceAndIndex(price, pricePoint, i)
			if len(pricePoint.Entries) == 0 && book.HighestLossPrice == price {
				if ok := iterator.Previous(); ok {
					book.HighestLossPrice = iterator.Key()
				} else {
					book.HighestLossPrice = 0
				}
			}
			iterator.Close()
			return true
		}
	}
	iterator.Close()
	return false
}

// Cancel stop entry order and return true if found
func (book *orderBook) cancelStopEntryOrder(order model.Order, events *[]model.Event) bool {
	price := order.StopPrice
	iterator := book.StopEntryOrders.Seek(price)
	// price is outside the bounds of the list
	if iterator == nil {
		return false
	}
	// price is in the range but does not exist in the list
	if iterator.Key() != price {
		iterator.Close()
		return false
	}
	pricePoint := iterator.Value()
	for i := 0; i < len(pricePoint.Entries); i++ {
		if pricePoint.Entries[i].ID == order.ID {
			ord := pricePoint.Entries[i]
			ord.SetStatus(model.OrderStatus_Cancelled)
			book.LastEventSeqID++
			*events = append(*events, model.NewOrderStatusEvent(book.LastEventSeqID, book.MarketID, ord.Type, ord.Side, ord.ID, ord.OwnerID, ord.Price, ord.Amount, ord.Funds, ord.Status, ord.FilledAmount, ord.UsedFunds))
			book.StopEntryOrders.removeEntryByPriceAndIndex(price, pricePoint, i)
			if len(pricePoint.Entries) == 0 && book.LowestEntryPrice == price {
				if ok := iterator.Next(); ok {
					book.LowestEntryPrice = iterator.Key()
				} else {
					book.LowestEntryPrice = 0
				}
			}
			iterator.Close()
			return true
		}
	}
	iterator.Close()
	return false
}

// Cancel a stop order based on a given order ID and set price
func (book *orderBook) cancelStopOrder(order model.Order, events *[]model.Event) {
	if order.Stop == model.StopLoss_None {
		return
	}
	if order.Stop == model.StopLoss_Loss && book.cancelStopLossOrder(order, events) {
		return
	}
	if order.Stop == model.StopLoss_Entry && book.cancelStopEntryOrder(order, events) {
		return
	}
	// if nothing was cancelled is possible the order was already activated and we should cancel it from the market
	order.Stop = model.StopLoss_None
	if book.cancelLimitOrder(order, events) {
		return
	}
	book.AppendErrorEvent(events, model.ErrorCode_CancelFailed, order)
}

// Cancel a limit order based on a given order ID and set price
func (book *orderBook) cancelLimitOrder(order model.Order, events *[]model.Event) bool {
	if order.Side == model.MarketSide_Buy {
		iterator := book.BuyEntries.Seek(order.Price)
		// price is outside the bounds of the list
		if iterator == nil {
			return false
		}
		// price is in the range but does not exist in the list
		if iterator.Key() != order.Price {
			iterator.Close()
			return false
		}
		pricePoint := iterator.Value()
		for i := 0; i < len(pricePoint.Entries); i++ {
			if pricePoint.Entries[i].ID == order.ID {
				ord := pricePoint.Entries[i]
				ord.SetStatus(model.OrderStatus_Cancelled)
				book.LastEventSeqID++
				*events = append(*events, model.NewOrderStatusEvent(book.LastEventSeqID, book.MarketID, ord.Type, ord.Side, ord.ID, ord.OwnerID, ord.Price, ord.Amount, ord.Funds, ord.Status, ord.FilledAmount, ord.UsedFunds))
				book.removeBuyBookEntry(ord.Price, pricePoint, i)
				// adjust highest bid
				if len(pricePoint.Entries) == 0 && book.HighestBid == ord.Price {
					book.closeBidIterator(iterator)
				} else {
					iterator.Close()
				}
				return true
			}
		}
		iterator.Close()
		return false
	}

	// cancel sell limit order

	iterator := book.SellEntries.Seek(order.Price)
	// price is outside the bounds of the list
	if iterator == nil {
		return false
	}
	// price is in the range but does not exist in the list
	if iterator.Key() != order.Price {
		iterator.Close()
		return false
	}
	pricePoint := iterator.Value()
	for i := 0; i < len(pricePoint.Entries); i++ {
		if pricePoint.Entries[i].ID == order.ID {
			ord := pricePoint.Entries[i]
			ord.SetStatus(model.OrderStatus_Cancelled)
			book.LastEventSeqID++
			*events = append(*events, model.NewOrderStatusEvent(book.LastEventSeqID, book.MarketID, ord.Type, ord.Side, ord.ID, ord.OwnerID, ord.Price, ord.Amount, ord.Funds, ord.Status, ord.FilledAmount, ord.UsedFunds))
			book.removeSellBookEntry(ord.Price, pricePoint, i)
			// adjust lowest ask
			if len(pricePoint.Entries) == 0 && book.LowestAsk == ord.Price {
				book.closeAskIterator(iterator)
			} else {
				iterator.Close()
			}
			return true
		}
	}
	iterator.Close()
	return false
}

// Append an error event and increases alst event seq id
func (book *orderBook) AppendErrorEvent(events *[]model.Event, code model.ErrorCode, order model.Order) {
	book.LastEventSeqID++
	*events = append(*events, model.NewErrorEvent(book.LastEventSeqID, book.MarketID, code, order.Type, order.Side, order.ID, order.OwnerID, order.Price, order.Amount, order.Funds))
}

// Generate cancel order event, add it to the list of events and increment the LastEventSeqID
func (book *orderBook) generateCancelOrderEvent(order model.Order, events *[]model.Event) {
	book.LastEventSeqID++
	*events = append(*events, model.NewOrderStatusEvent(
		book.LastEventSeqID,
		book.MarketID,
		order.Type,
		order.Side,
		order.ID,
		order.OwnerID,
		order.Price,
		order.Amount,
		order.Funds,
		model.OrderStatus_Cancelled,
		order.FilledAmount,
		order.UsedFunds,
	))
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
	} else {
		book.HighestBid = 0
	}
	iterator.Close()
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
	} else {
		book.LowestAsk = 0
	}
	iterator.Close()
}
