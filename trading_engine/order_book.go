package trading_engine

import (
	"github.com/ryszard/goskiplist/skiplist"
)

// OrderBook type
type OrderBook struct {
	PricePoints *skiplist.SkipList
	LowestAsk   float64
	HighestBid  float64
	MarketPrice float64
}

// NewOrderBook Creates a new empty order book for the trading engine
func NewOrderBook() OrderBook {
	pricePoints := NewPricePoints()
	return OrderBook{PricePoints: pricePoints, LowestAsk: 0, HighestBid: 0, MarketPrice: 0}
}

func (orderBook *OrderBook) addBookEntry(bookEntry *BookEntry) {
	if value, ok := orderBook.PricePoints.Get(bookEntry.Price); ok {
		pricePoint := value.(*PricePoint)
		if bookEntry.Order.Side == BUY {
			pricePoint.BuyBookEntries = append(pricePoint.BuyBookEntries, bookEntry)
		} else {
			pricePoint.SellBookEntries = append(pricePoint.SellBookEntries, bookEntry)
		}
	} else {
		buyBookEntries := []*BookEntry{}
		sellBookEntries := []*BookEntry{}
		if bookEntry.Order.Side == BUY {
			buyBookEntries = append(buyBookEntries, bookEntry)
		} else {
			sellBookEntries = append(sellBookEntries, bookEntry)
		}
		pricePoint := &PricePoint{BuyBookEntries: buyBookEntries, SellBookEntries: sellBookEntries}
		orderBook.PricePoints.Set(bookEntry.Price, pricePoint)
	}
}

func (orderBook *OrderBook) removeBookEntry(bookEntry *BookEntry) {
	if value, ok := orderBook.PricePoints.Get(bookEntry.Price); ok {
		pricePoint := value.(*PricePoint)
		if bookEntry.Order.Side == BUY {
			for i, buyEntry := range pricePoint.BuyBookEntries {
				if buyEntry == bookEntry {
					pricePoint.BuyBookEntries = append(pricePoint.BuyBookEntries[:i], pricePoint.BuyBookEntries[i+1:]...)
					break
				}
			}
		} else {
			for i, sellEntry := range pricePoint.SellBookEntries {
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

// Process a new received order and return a list of trades make
func (orderBook *OrderBook) Process(order Order) []Trade {
	var trades []Trade
	if order.Side == BUY {
		trades = orderBook.processLimitBuy(order)
	} else {
		trades = orderBook.processLimitSell(order)
	}
	return trades
}

func (orderBook *OrderBook) processLimitBuy(order Order) []Trade {
	bookEntry := NewBookEntry(order)
	var trades []Trade

	if orderBook.LowestAsk <= bookEntry.Price && orderBook.LowestAsk != 0 {
		iterator := orderBook.PricePoints.Seek(orderBook.LowestAsk)
		defer iterator.Close()

		// traverse orders to find a matching one based on the sell order list
		for iterator != nil && bookEntry.Price >= orderBook.LowestAsk {
			pricePoint := iterator.Value().(*PricePoint)
			for _, sellEntry := range pricePoint.SellBookEntries {
				// if we can fill the trade instantly then we add the trade and complete the order
				if sellEntry.Amount >= bookEntry.Amount {
					trades = append(trades, NewTrade(bookEntry.Order, sellEntry.Order, bookEntry.Amount, bookEntry.Price))
					sellEntry.Amount -= bookEntry.Amount
					if sellEntry.Amount == 0 {
						orderBook.removeBookEntry(sellEntry)
					}
					return trades
				}

				// if the sell order has a lower amount than what the buy order is then we fill only what we can from the sell order,
				// we complete the sell order and we move to the next order
				if sellEntry.Amount < bookEntry.Amount {
					trades = append(trades, NewTrade(bookEntry.Order, sellEntry.Order, sellEntry.Amount, sellEntry.Price))
					bookEntry.Amount -= sellEntry.Amount
					orderBook.removeBookEntry(sellEntry)
					continue
				}
			}

			if ok := iterator.Next(); ok {
				if len(iterator.Value().(*PricePoint).SellBookEntries) > 0 {
					orderBook.LowestAsk = iterator.Key().(float64)
				}
			} else {
				orderBook.LowestAsk = 0
				break
			}
		}
	}

	// if there are no more ordes just add the buy order to the list
	orderBook.addBookEntry(bookEntry)
	if orderBook.HighestBid < order.Price || orderBook.HighestBid == 0 {
		orderBook.HighestBid = order.Price
	}

	return trades
}

func (orderBook *OrderBook) processLimitSell(order Order) []Trade {
	bookEntry := NewBookEntry(order)
	var trades []Trade

	if orderBook.HighestBid >= bookEntry.Price && orderBook.HighestBid != 0 {
		iterator := orderBook.PricePoints.Seek(orderBook.HighestBid)
		defer iterator.Close()

		// traverse orders to find a matching one based on the sell order list
		for iterator != nil && bookEntry.Price <= orderBook.HighestBid {
			pricePoint := iterator.Value().(*PricePoint)
			for _, buyEntry := range pricePoint.BuyBookEntries {
				// if we can fill the trade instantly then we add the trade and complete the order
				if buyEntry.Amount >= bookEntry.Amount {
					trades = append(trades, NewTrade(bookEntry.Order, buyEntry.Order, bookEntry.Amount, bookEntry.Price))
					buyEntry.Amount -= bookEntry.Amount
					if buyEntry.Amount == 0 {
						orderBook.removeBookEntry(buyEntry)
					}
					return trades
				}

				// if the sell order has a lower amount than what the buy order is then we fill only what we can from the sell order,
				// we complete the sell order and we move to the next order
				if buyEntry.Amount < bookEntry.Amount {
					trades = append(trades, NewTrade(bookEntry.Order, buyEntry.Order, buyEntry.Amount, buyEntry.Price))
					bookEntry.Amount -= buyEntry.Amount
					orderBook.removeBookEntry(buyEntry)
					continue
				}
			}

			if ok := iterator.Previous(); ok {
				if len(iterator.Value().(*PricePoint).BuyBookEntries) > 0 {
					orderBook.HighestBid = iterator.Key().(float64)
				}
			} else {
				orderBook.HighestBid = 0
				break
			}
		}
	}

	// if there are no more ordes just add the buy order to the list
	orderBook.addBookEntry(bookEntry)
	if orderBook.LowestAsk > order.Price || orderBook.LowestAsk == 0 {
		orderBook.LowestAsk = order.Price
	}

	return trades
}
