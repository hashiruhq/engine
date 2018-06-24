package trading_engine

import (
	"github.com/ryszard/goskiplist/skiplist"
)

// OrderBook type
type OrderBook struct {
	PricePoints *skiplist.SkipList
	LowestAsk   uint64
	HighestBid  uint64
	MarketPrice uint64
}

// NewOrderBook Creates a new empty order book for the trading engine
func NewOrderBook() OrderBook {
	pricePoints := NewPricePoints()
	return OrderBook{PricePoints: pricePoints, LowestAsk: 0, HighestBid: 0, MarketPrice: 0}
}

// Process a new received order and return a list of trades make
func (orderBook *OrderBook) Process(order Order) []Trade {
	if order.Side == BUY {
		return orderBook.processLimitBuy(order)
	}
	return orderBook.processLimitSell(order)
}

func (orderBook *OrderBook) processLimitBuy(order Order) []Trade {
	var trades []Trade
	if orderBook.LowestAsk <= order.Price && orderBook.LowestAsk != 0 {
		iterator := orderBook.PricePoints.Seek(orderBook.LowestAsk)

		// traverse orders to find a matching one based on the sell order list
		if iterator != nil {
			defer iterator.Close()
			for order.Price >= orderBook.LowestAsk {
				pricePoint := iterator.Value().(*PricePoint)
				for index := 0; index < len(pricePoint.SellBookEntries); index++ {
					sellEntry := &pricePoint.SellBookEntries[index]
					// if we can fill the trade instantly then we add the trade and complete the order
					if sellEntry.Amount >= order.Amount {
						trades = append(trades, NewTrade(order.ID, sellEntry.Order.ID, order.Amount, sellEntry.Order.Price))
						sellEntry.Amount -= order.Amount
						if sellEntry.Amount == 0 {
							orderBook.removeBookEntry(sellEntry)
						}
						return trades
					}

					// if the sell order has a lower amount than what the buy order is then we fill only what we can from the sell order,
					// we complete the sell order and we move to the next order
					if sellEntry.Amount < order.Amount {
						trades = append(trades, NewTrade(order.ID, sellEntry.Order.ID, sellEntry.Amount, sellEntry.Price))
						order.Amount -= sellEntry.Amount
						index--
						orderBook.removeBookEntry(sellEntry)
						continue
					}
				}

				if ok := iterator.Next(); ok {
					if len(iterator.Value().(*PricePoint).SellBookEntries) > 0 {
						orderBook.LowestAsk = iterator.Key().(uint64)
					}
				} else {
					orderBook.LowestAsk = 0
					break
				}
			}
		}
	}

	// if there are no more ordes just add the buy order to the list
	orderBook.addBookEntry(BookEntry{Price: order.Price, Amount: order.Amount, Order: order})
	if orderBook.HighestBid < order.Price || orderBook.HighestBid == 0 {
		orderBook.HighestBid = order.Price
	}

	return trades
}

func (orderBook *OrderBook) processLimitSell(order Order) []Trade {
	var trades []Trade
	if orderBook.HighestBid >= order.Price && orderBook.HighestBid != 0 {
		iterator := orderBook.PricePoints.Seek(orderBook.HighestBid)

		// traverse orders to find a matching one based on the sell order list
		if iterator != nil {
			defer iterator.Close()
			for order.Price <= orderBook.HighestBid {
				pricePoint := iterator.Value().(*PricePoint)
				for index := 0; index < len(pricePoint.BuyBookEntries); index++ {
					buyEntry := &pricePoint.BuyBookEntries[index]
					// if we can fill the trade instantly then we add the trade and complete the order
					if buyEntry.Amount >= order.Amount {
						trades = append(trades, NewTrade(order.ID, buyEntry.Order.ID, order.Amount, buyEntry.Order.Price))
						buyEntry.Amount -= order.Amount
						if buyEntry.Amount == 0 {
							orderBook.removeBookEntry(buyEntry)
						}
						return trades
					}

					// if the sell order has a lower amount than what the buy order is then we fill only what we can from the sell order,
					// we complete the sell order and we move to the next order
					if buyEntry.Amount < order.Amount {
						trades = append(trades, NewTrade(order.ID, buyEntry.Order.ID, buyEntry.Amount, buyEntry.Price))
						order.Amount -= buyEntry.Amount
						orderBook.removeBookEntry(buyEntry)
						index--
						continue
					}
				}

				if ok := iterator.Previous(); ok {
					if len(iterator.Value().(*PricePoint).BuyBookEntries) > 0 {
						orderBook.HighestBid = iterator.Key().(uint64)
					}
				} else {
					orderBook.HighestBid = 0
					break
				}
			}
		}
	}

	// if there are no more ordes just add the buy order to the list
	orderBook.addBookEntry(BookEntry{Price: order.Price, Amount: order.Amount, Order: order})
	if orderBook.LowestAsk > order.Price || orderBook.LowestAsk == 0 {
		orderBook.LowestAsk = order.Price
	}

	return trades
}

// Cancel an order from the order book based on the order price and ID
func (orderBook *OrderBook) Cancel(id string, price uint64) error {
	// @todo implement this method
	return nil
}

// Add a new book entry in the order book
// If the price point already exists then the book entry is simply added at the end of the pricepoint
// If the price point does not exist yet it will be created
func (orderBook *OrderBook) addBookEntry(bookEntry BookEntry) {
	var pricePoint *PricePoint
	price := bookEntry.Price
	if value, ok := orderBook.PricePoints.Get(price); ok {
		pricePoint = value.(*PricePoint)
	} else {
		pricePoint = &PricePoint{
			BuyBookEntries:  []BookEntry{},
			SellBookEntries: []BookEntry{},
		}
		orderBook.PricePoints.Set(price, pricePoint)
	}

	if bookEntry.Order.Side == BUY {
		pricePoint.BuyBookEntries = append(pricePoint.BuyBookEntries, bookEntry)
	} else {
		pricePoint.SellBookEntries = append(pricePoint.SellBookEntries, bookEntry)
	}
}

// Remove a book entry from the order book
// The method will also remove the price point entry if the
func (orderBook *OrderBook) removeBookEntry(bookEntry *BookEntry) {
	if value, ok := orderBook.PricePoints.Get(bookEntry.Price); ok {
		pricePoint := value.(*PricePoint)
		if bookEntry.Order.Side == BUY {
			for i := range pricePoint.BuyBookEntries {
				buyEntry := &pricePoint.BuyBookEntries[i]
				if buyEntry == bookEntry {
					pricePoint.BuyBookEntries = append(pricePoint.BuyBookEntries[:i], pricePoint.BuyBookEntries[i+1:]...)
					break
				}
			}
		} else {
			for i := range pricePoint.SellBookEntries {
				sellEntry := &pricePoint.SellBookEntries[i]
				if bookEntry == sellEntry {
					pricePoint.SellBookEntries = append(pricePoint.SellBookEntries[:i], pricePoint.SellBookEntries[i+1:]...)
					break
				}
			}
		}
		if len(pricePoint.BuyBookEntries) == 0 && len(pricePoint.SellBookEntries) == 0 {
			orderBook.PricePoints.Delete(bookEntry.Price)
		}
	}
}
