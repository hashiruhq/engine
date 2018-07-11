package trading_engine

// import "github.com/ryszard/goskiplist/skiplist"
// import "trading_engine/skiplist"

// PricePoint holds the records for a particular price
type PricePoint struct {
	SellBookEntries []BookEntry
	BuyBookEntries  []BookEntry
}

// NewPricePoints creates a new skiplist in which to hold all price points
func NewPricePoints() *SkipList {
	return NewCustomMap(func(l, r uint64) bool {
		return l < r
	})
}
