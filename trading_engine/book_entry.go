package trading_engine

import (
	"gitlab.com/around25/products/matching-engine/conv"

	"github.com/francoispqt/gojay"
)

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

// FromJSON loads an entry from a byte array
func (entry *BookEntry) FromJSON(msg []byte) error {
	return gojay.Unsafe.Unmarshal(msg, entry)
}

// ToJSON converts an entry to a byte string
func (entry *BookEntry) ToJSON() ([]byte, error) {
	return gojay.Marshal(entry)
}

// UnmarshalJSONObject implement gojay.UnmarshalerJSONObject interface
func (entry *BookEntry) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "order":
		return dec.Object(&entry.Order)
	case "price":
		var amount string
		dec.String(&amount)
		entry.Price = conv.ToUnits(amount, PricePrecision)
	case "amount":
		var amount string
		dec.String(&amount)
		entry.Amount = conv.ToUnits(amount, PricePrecision)
	}
	return nil
}

// NKeys implements gojay.UnmarshalerJSONObject interface and returns the number of keys to parse
func (entry BookEntry) NKeys() int {
	return 3
}

// MarshalJSONObject implement gojay.MarshalerJSONObject interface
func (entry BookEntry) MarshalJSONObject(enc *gojay.Encoder) {
	enc.ObjectKey("order", entry.Order)
	enc.StringKey("price", conv.FromUnits(entry.Price, PricePrecision))
	enc.StringKey("amount", conv.FromUnits(entry.Amount, AmountPrecision))
}

// IsNil checks if the order is empty
func (entry BookEntry) IsNil() bool {
	return entry.Order.IsNil()
}

// BookEntries defines convert operations on book entries
type BookEntries []BookEntry

// implement UnmarshalerJSONArray
func (bookEntries *BookEntries) UnmarshalJSONArray(dec *gojay.Decoder) error {
	entry := BookEntry{}
	if err := dec.Object(&entry); err != nil {
		return err
	}
	*bookEntries = append(*bookEntries, entry)
	return nil
}

// implement MarshalerJSONArray
func (bookEntries BookEntries) MarshalJSONArray(enc *gojay.Encoder) {
	for _, entry := range bookEntries {
		enc.Object(entry)
	}
}
func (bookEntries BookEntries) IsNil() bool {
	return len(bookEntries) == 0
}
