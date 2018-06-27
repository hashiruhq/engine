package trading_engine

import (
	"trading_engine/conv"

	"github.com/francoispqt/gojay"
)

// BUY value means the user wants to buy from the market
const BUY = 1

// SELL value means the user wants to sell to the market
const SELL = 2

//LimitOrder allows the trader to start an order where the transaction will be completed
// if the market price is at or better than the set price
const LimitOrder = 1

// MarketOrder completes the trade at the current market price
const MarketOrder = 2

// StopLossOrder @todo completes the trade until it gets to a price
const StopLossOrder = 3

// PricePrecision - The maximum precision of the price for an item in the market
// Max Price: 184467440737.09551615
const PricePrecision = 8

// AmountPrecision - The maximum precision of the amount for an item in the market
// Max Amount: 184467440737.09551615
const AmountPrecision = 8

// FundsPrecision - the precision for the provided funds
const FundsPrecision = 8

// Order allows the trader to start an order where the transaction will be completed
// if the market price is at or better than the set price
type Order struct {
	// Optional:
	// Amount of coins to buy/sell with the order
	// - The amount must be greater than the base_min_amount for the product and no larger than the base_max_amount.
	Amount uint64

	// Optional:
	// The price to pay for one unit in the market
	// - The price must be specified in quote_increment product units.
	// - The quote increment is the smallest unit of price. For the BTC-USD product,
	//   the quote increment is 0.01 or 1 penny. Prices less than 1 penny will not be accepted,
	//   and no fractional penny prices will be accepted. Not required for market orders.
	Price uint64

	//******************************************
	// Common fields
	//******************************************
	// The id of the order
	ID string

	// User Order ID: An id defined by the user to identity the order
	// UserOrderId string

	// The type of the order: 1=limit 2=market
	Type int8
	// Category int8 // deprecated by the Type field
	// The side of the market: 1=buy 2=sell
	Side int8
	// Base and Quote symbols as lowercase 3-4 letter words. Ex: btc, usd, eth
	BaseCurrency  string
	QuoteCurrency string
	// FUTURE FIELD
	// Prevent self trade
	//
	// Self-trading is not allowed on the trading engine. Two orders from the same user will not fill one another.
	// When placing an order, you can specify the self-trade prevention behavior.
	//
	// DECREMENT AND CANCEL
	// The default behavior is decrement and cancel. When two orders from the same user cross, the smaller order
	// will be canceled and the larger order size will be decremented by the smaller order size. If the two orders
	// are the same size, both will be canceled.
	//
	// CANCEL OLDEST
	// Cancel the older (resting) order in full. The new order continues to execute.
	// CANCEL NEWEST
	// Cancel the newer (taking) order in full. The old resting order remains on the order book.

	// CANCEL BOTH
	// Immediately cancel both orders.

	// NOTES FOR MARKET ORDERS
	// When a market order using dc self-trade prevention encounters an open limit order, the behavior depends on
	// which fields for the market order message were specified. If funds and size are specified for a buy order,
	// then size for the market order will be decremented internally within the matching engine and funds will
	// remain unchanged. The intent is to offset your target size without limiting your buying power. If size is
	// not specified, then funds will be decremented. For a market sell, the size will be decremented when
	// encountering existing limit orders.
	// PreventSelfTrade string

	// Stop flag. Requires `StopPrice`` to be defined.
	// Stop orders become active and wait to trigger based on the movement of the last trade price.
	// THere are 2 types of stop orders: 1=loss 2=entry
	// - Stop loss triggers when the last trade price changes to a value at or below the `StopPrice`.
	// - Stop entry triggers when the last trade price changes to a value at or above the `StopPrice`.
	// - Note that when triggered, stop orders execute as either market or limit orders, depending on the type.
	Stop int8
	// Sets trigger price for stop order. Only if stop is defined.
	StopPrice uint64

	//*****************************************
	// Market Order Fields
	// - Requires the Amount field from above
	// - At lease one of Amount or Funds fields should be set
	// - Funds limit how much your quote currency account balance is used and
	//   Amount limits the amount of coins that will be transacted
	// - Market orders are always considered takers and should always fill immediately
	//*****************************************

	// Maximum total funds to use for the order
	// - The funds field is optionally used for market orders. When specified it indicates how much of the product
	//   quote currency to buy or sell. For example, a market buy for BTC-USD with funds specified as 150.00 will
	//   spend 150 USD to buy BTC (including any fees). If the funds field is not specified for a market buy order,
	//   size must be specified and the enting will use available funds in your account to buy bitcoin.
	// - A market sell order can also specify the funds. If funds is specified, it will limit the sell to the amount
	//   of funds specified. You can use funds with sell orders to limit the amount of quote currency funds received.
	Funds uint64

	//*****************************************
	// Limit Order
	// - Requires the Amount field from above
	// - Requires the Price field from above
	//*****************************************

	// FUTURE PROPERTY
	//
	// TimeInForce string // time in force. GTC, GTT, IOC, or FOK (default is GTC)

	// FUTURE PROPERTY
	//
	// CancelAfter uint   // cancel after min, hour, day. Requires time_in_force to be GTT

	// FUTURE PROPERTY
	// The post-only flag indicates that the order should only make liquidity. If any part of the
	// order results in taking liquidity, the order will be rejected and no part of it will execute.
	//
	// PostOnly    string // post only. Invalid when time_in_force is IOC or FOK
}

// NewOrder create a new order
func NewOrder(id string, price, amount uint64, side int8, category int8) Order {
	return Order{ID: id, Price: price, Amount: amount, Side: side, Type: category}
}

// LessThan implementes the skiplist interface
func (order Order) LessThan(other Order) bool {
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
	case "type":
		return dec.Int8(&order.Type)
	case "price":
		var price float64
		err := dec.Float(&price)
		if err != nil {
			return err
		}
		order.Price = conv.ToUnits(price, PricePrecision)
		return nil
	case "amount":
		var amount float64
		err := dec.Float(&amount)
		if err != nil {
			return err
		}
		order.Amount = conv.ToUnits(amount, AmountPrecision)
		return nil
	}
	return nil
}

// NKeys implements gojay.UnmarshalerJSONObject interface and returns the number of keys to parse
func (order Order) NKeys() int {
	return 5
}

// MarshalJSONObject implement gojay.MarshalerJSONObject interface
func (order Order) MarshalJSONObject(enc *gojay.Encoder) {
	enc.StringKey("id", order.ID)
	enc.IntKey("side", int(order.Side))
	enc.IntKey("type", int(order.Type))
	enc.FloatKey("price", conv.FromUnits(order.Price, PricePrecision))
	enc.FloatKey("amount", conv.FromUnits(order.Amount, AmountPrecision))
}

// IsNil checks if the order is empty
func (order *Order) IsNil() bool {
	return order == nil
}
