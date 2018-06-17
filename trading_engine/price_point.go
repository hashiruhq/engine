package trading_engine

import "github.com/ryszard/goskiplist/skiplist"

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
