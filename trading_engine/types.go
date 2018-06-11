package trading_engine

import (
	"sync"
	"time"

	"github.com/ryszard/goskiplist/skiplist"
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

// LessThan implementes the skiplist interface
func (order *Order) LessThan(other *Order) bool {
	return order.Price < other.Price
}

// PricePoint holds the records for a particular price
type PricePoint struct {
	SellBookEntries []*BookEntry
	BuyBookEntries  []*BookEntry
}

// NewPricePoints creates a new skiplist in which to hold all price points
func NewPricePoints() *skiplist.SkipList {
	return skiplist.NewCustomMap(func(l, r interface{}) bool {
		return l.(float64) < r.(float64)
	})
}

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

// BookEntry keeps minimal information about the order in the OrderBook
type BookEntry struct {
	Price  float64
	Amount float64
	Order  *Order
}

// NewBookEntry Creates a new book entry based on an Order to be used by the order book
func NewBookEntry(order *Order) *BookEntry {
	return &BookEntry{Price: order.Price, Amount: order.Amount, Order: order}
}

// OrderBook type
type OrderBook struct {
	mtx         sync.Mutex
	PricePoints *skiplist.SkipList
	// BuyOrders   []*BookEntry
	// SellOrders  []*BookEntry
	LowestAsk   float64
	HighestBid  float64
	MarketPrice float64
}

// NewOrderBook Creates a new empty order book for the trading engine
func NewOrderBook() *OrderBook {
	pricePoints := NewPricePoints()
	// buyOrders := make([]*BookEntry, 0, 100000)
	// sellOrders := make([]*BookEntry, 0, 100000)
	return &OrderBook{PricePoints: pricePoints, LowestAsk: 0, HighestBid: 0, MarketPrice: 0}
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
