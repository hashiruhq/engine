package trading_engine

// OrderBook interface
// Defines what constitudes an order book and how we can interact with it
type OrderBook interface {
	Process(Order) []Trade
	Cancel(id string) *BookEntry
	GetHighestBid() uint64
	GetLowestAsk() uint64
	Load(Market) error
	Backup() Market
	GetMarket() *SkipList
}

type orderBook struct {
	PricePoints *SkipList
	LowestAsk   uint64
	HighestBid  uint64
	OpenOrders  map[string]uint64
}

// NewOrderBook Creates a new empty order book for the trading engine
func NewOrderBook() OrderBook {
	pricePoints := NewPricePoints()
	return &orderBook{
		PricePoints: pricePoints,
		LowestAsk:   0,
		HighestBid:  0,
		OpenOrders:  make(map[string]uint64),
	}
}

func (book *orderBook) AddOpenOrder(id string, price uint64) {
	book.OpenOrders[id] = price
}

func (book *orderBook) RemoveOpenOrder(id string) {
	delete(book.OpenOrders, id)
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
func (book orderBook) GetMarket() *SkipList {
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
func (book *orderBook) Cancel(id string) *BookEntry {
	// @todo implement this method
	if price, ok := book.OpenOrders[id]; ok {
		if value, ok := book.PricePoints.Get(price); ok {
			pricePoint := value
			for i := 0; i < len(pricePoint.BuyBookEntries); i++ {
				if pricePoint.BuyBookEntries[i].Order.ID == id {
					bookEntry := pricePoint.BuyBookEntries[i]
					book.removeBuyBookEntry(&bookEntry, pricePoint, i)
					return &bookEntry
				}
			}
			for i := 0; i < len(pricePoint.SellBookEntries); i++ {
				if pricePoint.SellBookEntries[i].Order.ID == id {
					bookEntry := pricePoint.SellBookEntries[i]
					book.removeSellBookEntry(&bookEntry, pricePoint, i)
					return &bookEntry
				}
			}
		}
	}
	return nil
}

// Add a new book entry in the order book
// If the price point already exists then the book entry is simply added at the end of the pricepoint
// If the price point does not exist yet it will be created

func (book *orderBook) addBuyBookEntry(bookEntry BookEntry) {
	price := bookEntry.Price
	if value, ok := book.PricePoints.Get(price); ok {
		pricePoint := value
		pricePoint.BuyBookEntries = append(pricePoint.BuyBookEntries, bookEntry)
		book.AddOpenOrder(bookEntry.Order.ID, price)
		return
	}
	pricePoint := &PricePoint{
		BuyBookEntries:  []BookEntry{bookEntry},
		SellBookEntries: []BookEntry{},
	}
	book.PricePoints.Set(price, pricePoint)
	book.AddOpenOrder(bookEntry.Order.ID, price)
}

func (book *orderBook) addSellBookEntry(bookEntry BookEntry) {
	price := bookEntry.Price
	if value, ok := book.PricePoints.Get(price); ok {
		pricePoint := value
		pricePoint.SellBookEntries = append(pricePoint.SellBookEntries, bookEntry)
		book.AddOpenOrder(bookEntry.Order.ID, price)
		return
	}
	pricePoint := &PricePoint{
		BuyBookEntries:  []BookEntry{},
		SellBookEntries: []BookEntry{bookEntry},
	}
	book.PricePoints.Set(price, pricePoint)
	book.AddOpenOrder(bookEntry.Order.ID, price)
}

// Remove a book entry from the order book
// The method will also remove the price point entry if both book entry lists are empty

func (book *orderBook) removeBuyBookEntry(bookEntry *BookEntry, pricePoint *PricePoint, index int) {
	book.RemoveOpenOrder(bookEntry.Order.ID)
	pricePoint.BuyBookEntries = append(pricePoint.BuyBookEntries[:index], pricePoint.BuyBookEntries[index+1:]...)
	if len(pricePoint.BuyBookEntries) == 0 && len(pricePoint.SellBookEntries) == 0 {
		book.PricePoints.Delete(bookEntry.Price)
	}
}

func (book *orderBook) removeSellBookEntry(bookEntry *BookEntry, pricePoint *PricePoint, index int) {
	book.RemoveOpenOrder(bookEntry.Order.ID)
	pricePoint.SellBookEntries = append(pricePoint.SellBookEntries[:index], pricePoint.SellBookEntries[index+1:]...)
	if len(pricePoint.BuyBookEntries) == 0 && len(pricePoint.SellBookEntries) == 0 {
		book.PricePoints.Delete(bookEntry.Price)
	}
}
