package trading_engine

// Trade represents a completed trade between two orders
type Trade struct {
	TakerOrderID string  `json:"taker_order_id"`
	MakerOrderID string  `json:"maker_order_id"`
	Amount       float64 `json:"amount"`
	Price        float64 `json:"price"`
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
