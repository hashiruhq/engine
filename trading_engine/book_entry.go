package trading_engine

// BookEntry keeps minimal information about the order in the OrderBook
type BookEntry struct {
	Price  uint64
	Amount uint64
	Order  Order
}

// NewBookEntry Creates a new book entry based on an Order to be used by the order book
func NewBookEntry(order Order) BookEntry {
	return BookEntry{Price: order.Price, Amount: order.Amount, Order: order}
}
