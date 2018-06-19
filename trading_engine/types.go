package trading_engine

const (
	// BUY value means the user wants to buy from the market
	BUY = 1
	// SELL value means the user wants to sell to the market
	SELL = 2
)

const (
	//LIMIT_ORDER allows the trader to start an order where the transaction will be completed
	// if the market price is at or better than the set price
	LIMIT_ORDER = 1
	// MARKET_ORDER completes the trade at the current market price
	MARKET_ORDER = 2
	// STOP_LOSS_ORDER @todo completes the trade until it gets to a price
	STOP_LOSS_ORDER = 3
)

// // IOrder is defines what operations can be performed on an Order
// type IOrder interface {
// 	ComparePrice() int
// 	Subtract(order IOrder) IOrder
// }
