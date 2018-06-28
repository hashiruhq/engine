package trading_engine

import "testing"

func TestPricePointNew(t *testing.T) {
	pricePoints := NewPricePoints()
	order := NewOrder("TEST_1", uint64(100000000), uint64(12000000000), 1, 1)
	bookEntry := NewBookEntry(order)

	pricePoints.Set(bookEntry.Price, &PricePoint{
		BuyBookEntries: []BookEntry{bookEntry},
	})

	len := pricePoints.Len()
	if len != 1 {
		t.Errorf("Price points length invalid: got %v, want %v", len, 1)
	}

	value, ok := pricePoints.Get(uint64(100000000))
	if !ok {
		t.Errorf("Unable to get price point from list: got %v, want %v", ok, true)
	} else {
		point, _ := value.(*PricePoint)
		entry := point.BuyBookEntries[0]
		if entry.Price != 100000000 {
			t.Errorf("Book entry price invalid from price point: got %v, want %v", entry.Price, 100000000)
		}
	}

}
