package engine

import (
	"gitlab.com/around25/products/matching-engine/conv"

	"github.com/francoispqt/gojay"
)

// Market contains a user friendly representation of the order book
type Market struct {
	Topic      string
	Partition  int32
	Offset     int64
	LowestAsk  uint64
	HighestBid uint64
	BuyOrders  BookEntries
	SellOrders BookEntries
}

// FromJSON loads a market from a byte array
func (market *Market) FromJSON(msg []byte) error {
	return gojay.Unsafe.Unmarshal(msg, market)
}

// ToJSON converts a market to a byte string
func (market *Market) ToJSON() ([]byte, error) {
	return gojay.Marshal(market)
}

// UnmarshalJSONObject implement gojay.UnmarshalerJSONObject interface
func (market *Market) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "topic":
		return dec.String(&market.Topic)
	case "partition":
		return dec.Int32(&market.Partition)
	case "offset":
		return dec.Int64(&market.Offset)
	case "lowest_ask":
		var amount string
		dec.String(&amount)
		market.LowestAsk = conv.ToUnits(amount, PricePrecision)
	case "highest_bid":
		var amount string
		dec.String(&amount)
		market.HighestBid = conv.ToUnits(amount, PricePrecision)
	case "buy_orders":
		return dec.Array(&market.BuyOrders)
	case "sell_orders":
		return dec.Array(&market.SellOrders)
	}
	return nil
}

// NKeys implements gojay.UnmarshalerJSONObject interface and returns the number of keys to parse
func (market Market) NKeys() int {
	return 7
}

// MarshalJSONObject implement gojay.MarshalerJSONObject interface
func (market Market) MarshalJSONObject(enc *gojay.Encoder) {
	enc.StringKey("topic", market.Topic)
	enc.Uint64Key("partition", uint64(market.Partition))
	enc.Uint64Key("offset", uint64(market.Offset))
	enc.StringKey("lowest_ask", conv.FromUnits(market.LowestAsk, PricePrecision))
	enc.StringKey("highest_bid", conv.FromUnits(market.HighestBid, PricePrecision))
	enc.ArrayKey("buy_orders", market.BuyOrders)
	enc.ArrayKey("sell_orders", market.SellOrders)
}

// IsNil checks if the market is empty
func (market *Market) IsNil() bool {
	return market.Topic == "" && market.Partition == 0 && market.Offset == 0 && len(market.BuyOrders) == 0 && len(market.SellOrders) == 0
}
