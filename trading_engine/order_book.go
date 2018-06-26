package trading_engine

import (
	"github.com/ryszard/goskiplist/skiplist"
)

// OrderBook interface
// Defines what constitudes an order book and how we can interact with it
type OrderBook interface {
	Process(Order) []Trade
	Cancel(id string) error
	GetHighestBid() uint64
	GetLowestAsk() uint64
	GetMarket() *skiplist.SkipList
}

type orderBook struct {
	PricePoints *skiplist.SkipList
	LowestAsk   uint64
	HighestBid  uint64
	MarketPrice uint64
}

// NewOrderBook Creates a new empty order book for the trading engine
func NewOrderBook() OrderBook {
	pricePoints := NewPricePoints()
	return &orderBook{
		PricePoints: pricePoints,
		LowestAsk:   0,
		HighestBid:  0,
		MarketPrice: 0,
	}
}

// GetHighestBid returns the highest bid of the current market
func (book orderBook) GetHighestBid() uint64 {
	return book.HighestBid
}

// GetLowestAsk returns the lowest ask of the current market
func (book orderBook) GetLowestAsk() uint64 {
	return book.LowestAsk
}

// GetMarket returns the list of price points for the market
func (book orderBook) GetMarket() *skiplist.SkipList {
	return book.PricePoints
}

// Process a new received order and return a list of trades make
func (book *orderBook) Process(order Order) []Trade {
	if order.Side == BUY {
		return book.processLimitBuy(order)
	}
	return book.processLimitSell(order)
}

// Cancel an order from the order book based on the order price and ID
func (book *orderBook) Cancel(id string) error {
	// @todo implement this method
	return nil
}

// Add a new book entry in the order book
// If the price point already exists then the book entry is simply added at the end of the pricepoint
// If the price point does not exist yet it will be created
func (book *orderBook) addBookEntry(bookEntry BookEntry) {
	var pricePoint *PricePoint
	price := bookEntry.Price
	if value, ok := book.PricePoints.Get(price); ok {
		pricePoint = value.(*PricePoint)
	} else {
		pricePoint = &PricePoint{
			BuyBookEntries:  []BookEntry{},
			SellBookEntries: []BookEntry{},
		}
		book.PricePoints.Set(price, pricePoint)
	}

	if bookEntry.Order.Side == BUY {
		pricePoint.BuyBookEntries = append(pricePoint.BuyBookEntries, bookEntry)
	} else {
		pricePoint.SellBookEntries = append(pricePoint.SellBookEntries, bookEntry)
	}
}

// Remove a book entry from the order book
// The method will also remove the price point entry if the
func (book *orderBook) removeBookEntry(bookEntry *BookEntry) {
	if value, ok := book.PricePoints.Get(bookEntry.Price); ok {
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
			book.PricePoints.Delete(bookEntry.Price)
		}
	}
}
