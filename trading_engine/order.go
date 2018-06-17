package trading_engine

// Order allows the trader to start an order where the transaction will be completed
// if the market price is at or better than the set price
type Order struct {
	Price    float64 `json:"price"`
	Amount   float64 `json:"amount"`
	Side     int     `json:"side"`
	Category int     `json:"category"`
	ID       string  `json:"id"`
	// Symbol   string        `json:"symbol"` // @todo: Not sure what this is used for
	// TraderID string        `json:"trader_id"` // @todo: we may not need this for the trading engine, instead we can use just the order id
}

// NewOrder create a new order
func NewOrder(id string, price, amount float64, side int, category int) Order {
	return Order{ID: id, Price: price, Amount: amount, Side: side, Category: category}
}

// LessThan implementes the skiplist interface
func (order *Order) LessThan(other *Order) bool {
	return order.Price < other.Price
}
