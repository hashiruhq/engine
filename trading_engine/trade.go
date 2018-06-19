package trading_engine

import (
	"github.com/francoispqt/gojay"
)

// Trade represents a completed trade between two orders
type Trade struct {
	TakerOrderID string
	MakerOrderID string
	Amount       float64
	Price        float64
	// Date         time.Time `json:"created_at"`
}

// NewTrade Creates a new trade between the taker order and the maker order
func NewTrade(takerOrder Order, makerOrder Order, amount float64, price float64) Trade {
	return Trade{
		TakerOrderID: takerOrder.ID,
		MakerOrderID: makerOrder.ID,
		Amount:       amount,
		Price:        price,
		// Date:         time.Now(), // do not generate the timestamp in the trading engine to increase performance
	}
}

// FromJSON loads a trade from a byte array
func (trade *Trade) FromJSON(msg []byte) error {
	return gojay.Unsafe.Unmarshal(msg, trade)
}

// ToJSON Converts the trade into a JSON byte array
func (trade *Trade) ToJSON() ([]byte, error) {
	return gojay.MarshalJSONObject(trade)
}

// UnmarshalJSONObject implement gojay.UnmarshalerJSONObject interface
func (trade *Trade) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "taker_order_id":
		return dec.String(&trade.TakerOrderID)
	case "maker_order_id":
		return dec.String(&trade.MakerOrderID)
	case "price":
		return dec.Float(&trade.Price)
	case "amount":
		return dec.Float(&trade.Amount)
	}
	return nil
}

// MarshalJSONObject implement gojay.MarshalerJSONObject interface
func (trade *Trade) MarshalJSONObject(enc *gojay.Encoder) {
	enc.StringKey("taker_order_id", trade.TakerOrderID)
	enc.StringKey("maker_order_id", trade.MakerOrderID)
	enc.FloatKey("price", trade.Price)
	enc.FloatKey("amount", trade.Amount)
}

// IsNil checks if the trade is empty
func (trade *Trade) IsNil() bool {
	return trade == nil
}

// NKeys implements gojay.UnmarshalerJSONObject interface and returns the number of keys to parse
func (trade *Trade) NKeys() int {
	return 4
}
