package engine

import (
	"gitlab.com/around25/products/matching-engine/model"
)

// import "github.com/ryszard/goskiplist/skiplist"

// PricePoint holds the records for a particular price
type PricePoint struct {
	Entries []model.Order
}

// NewPricePoints creates a new skiplist in which to hold all price points
func NewPricePoints() *SkipList {
	return NewCustomMap(cmp_func)
}

func cmp_func(l, r uint64) bool {
	return l < r
}
