package trading_engine

import (
	"sync"
	"time"
)

// import "github.com/ryszard/goskiplist/skiplist"

// OrderSide will give us the type of the order
type OrderSide bool

// OrderCategory will store the type of the order: Limit, Market, etc
type OrderCategory int

const (
	// BUY value means the user wants to buy from the market
	BUY = false
	// SELL value means the user wants to sell to the market
	SELL = true
)

const (
	//LIMIT_ORDER allows the trader to start an order where the transaction will be completed
	// if the market price is at or better than the set price
	LIMIT_ORDER = 0
	// MARKET_ORDER completes the trade at the current market price
	MARKET_ORDER = 1
	// STOP_LOSS_ORDER @todo completes the trade until it gets to a price
	STOP_LOSS_ORDER = 2
)

// IOrder is defines what operations can be performed on an Order
type IOrder interface {
	ComparePrice() int
	Subtract(order IOrder) IOrder
}

// Order allows the trader to start an order where the transaction will be completed
// if the market price is at or better than the set price
type Order struct {
	ID       string        `json:"id"`
	Side     OrderSide     `json:"side"`
	Category OrderCategory `json:"category"`
	Amount   float64       `json:"amount"`
	Price    float64       `json:"price"`
	// Symbol   string        `json:"symbol"` // @todo: Not sure what this is used for
	// TraderID string        `json:"trader_id"` // @todo: we may not need this for the trading engine, instead we can use just the order id
}

// // LessThan implementes the skiplist interface
// func (order *Order) LessThan(other *Order) bool {
// 	return order.Price < other.Price
// }

// Trade represents a completed trade between two orders
type Trade struct {
	TakerOrderID string    `json:"taker_order_id"`
	MakerOrderID string    `json:"maker_order_id"`
	Amount       float64   `json:"amount"`
	Price        float64   `json:"price"`
	Date         time.Time `json:"created_at"`
}

// NewTrade Creates a new trade between the taker order and the maker order
func NewTrade(takerOrder *Order, makerOrder *Order, amount float64, price float64) *Trade {
	return &Trade{
		TakerOrderID: takerOrder.ID,
		MakerOrderID: makerOrder.ID,
		Amount:       amount,
		Price:        price,
		Date:         time.Now(),
	}
}

// // BookEntry keeps minimal information about the order in the OrderBook
// type BookEntry struct {
// 	Price  float64
// 	Amount float64
// 	Order  *Order
// }

// // NewBookEntry Creates a new book entry based on an Order to be used by the order book
// func NewBookEntry(order *Order) *BookEntry {
// 	return &BookEntry{Price: order.Price, Amount: order.Amount, Order: order}
// }

// OrderBook type
type OrderBook struct {
	mtx        sync.Mutex
	BuyOrders  []*Order
	SellOrders []*Order
}

// NewOrderBook Creates a new empty order book for the trading engine
func NewOrderBook() *OrderBook {
	buyOrders := make([]*Order, 0, 100000)
	sellOrders := make([]*Order, 0, 100000)
	return &OrderBook{BuyOrders: buyOrders, SellOrders: sellOrders}
}

// TradingEngine contains the current order book and information about the service since it was created
type TradingEngine struct {
	OrderBook       *OrderBook
	mtx             sync.Mutex
	Symbol          string
	LastRestart     string
	OrdersCompleted int64
	TradesCompleted int64
}

// NewTradingEngine creates a new trading engine that contains an empty order book and can start receving requests
func NewTradingEngine() *TradingEngine {
	orderBook := NewOrderBook()
	return &TradingEngine{OrderBook: orderBook, OrdersCompleted: 0, TradesCompleted: 0}
}
