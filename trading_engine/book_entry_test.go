package trading_engine

import "testing"

func TestBookEntry(t *testing.T) {
	// new order for 1.0 and 120 units
	order := NewOrder("TEST_1", 100000000, 12000000000, 1, 1)
	bookEntry := NewBookEntry(order)
	if bookEntry.Amount != 12000000000 {
		t.Errorf("Book entry amount incorrect, got: %d, want: %d.", bookEntry.Amount, 12000000000)
	}
	if bookEntry.Price != 100000000 {
		t.Errorf("Book entry price incorrect, got: %d, want: %d.", bookEntry.Amount, 100000000)
	}
	if bookEntry.Order != order {
		t.Errorf("Book entry order does not match, got: %v, want: %v.", bookEntry.Order, order)
	}
}
