package trading_engine

import "github.com/francoispqt/gojay"

// Order allows the trader to start an order where the transaction will be completed
// if the market price is at or better than the set price
type Order struct {
	Price    float64
	Amount   float64
	Side     int8
	Category int8
	ID       string
	// Symbol   string        `json:"symbol"` // @todo: Not sure what this is used for
	// TraderID string        `json:"trader_id"` // @todo: we may not need this for the trading engine, instead we can use just the order id
}

// NewOrder create a new order
func NewOrder(id string, price, amount float64, side int8, category int8) Order {
	return Order{ID: id, Price: price, Amount: amount, Side: side, Category: category}
}

// LessThan implementes the skiplist interface
func (order *Order) LessThan(other *Order) bool {
	return order.Price < other.Price
}

// FromJSON loads an order from a byte array
func (order *Order) FromJSON(msg []byte) error {
	return gojay.Unsafe.Unmarshal(msg, order)
}

// UnmarshalJSONObject implement gojay.UnmarshalerJSONObject interface
func (order *Order) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "id":
		return dec.String(&order.ID)
	case "side":
		return dec.Int8(&order.Side)
	case "category":
		return dec.Int8(&order.Category)
	case "price":
		return dec.Float(&order.Price)
	case "amount":
		return dec.Float(&order.Amount)
	}
	return nil
}

// NKeys implements gojay.UnmarshalerJSONObject interface and returns the number of keys to parse
func (order *Order) NKeys() int {
	return 5
}
