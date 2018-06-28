package trading_engine

import "testing"

func TestOrderNew(t *testing.T) {
	// new order for 1.0 and 120 units
	order := NewOrder("TEST_1", 100000000, 12000000000, 1, 2)
	if order.Amount != 12000000000 {
		t.Errorf("Order amount incorrect, got: %d, want: %d.", order.Amount, 12000000000)
	}
	if order.Price != 100000000 {
		t.Errorf("Order price incorrect, got: %d, want: %d.", order.Price, 100000000)
	}

	if order.Side != 1 {
		t.Errorf("Order price incorrect, got: %d, want: %d.", order.Side, 1)
	}

	if order.Type != 2 {
		t.Errorf("Order price incorrect, got: %d, want: %d.", order.Type, 2)
	}
}

func TestOrderLessThan(t *testing.T) {
	order1 := NewOrder("TEST_1", 100000000, 12000000000, 1, 2)
	order2 := NewOrder("TEST_2", 110000000, 12000000000, 1, 2)

	less := order1.LessThan(order2)
	if !less {
		t.Errorf("Abnormal order comparison, got: %v, want: %v.", less, true)
	}
}

func TestOrderFromJson(t *testing.T) {
	var order Order
	json := []byte(`{"base": "sym", "quote": "tst", "id":"TST_1", "price": 1312213.00010201, "amount": 8483828.29993942, "side": 1, "type": 1}`)

	order.FromJSON(json)

	if order.Price != 131221300010201 {
		t.Errorf("Failed to load order price from json, got: %v, want: %v.", order.Amount, 131221300010201)
	}

	if order.Amount != 848382829993942 {
		t.Errorf("Failed to load order amount from json, got: %v, want: %v.", order.Amount, 848382829993942)
	}

	if order.ID != "TST_1" {
		t.Errorf("Failed to load order id from json, got: %q, want: %q.", order.ID, "TST_1")
	}

	if order.Side != 1 {
		t.Errorf("Failed to load order side from json, got: %v, want: %v.", order.Side, 1)
	}

	if order.Type != 1 {
		t.Errorf("Failed to load order type from json, got: %v, want: %v.", order.Type, 1)
	}

	if order.BaseCurrency != "sym" {
		t.Errorf("Failed to load order base from json, got: %q, want: %q.", order.BaseCurrency, "sym")
	}

	if order.QuoteCurrency != "tst" {
		t.Errorf("Failed to load order quote from json, got: %q, want: %q.", order.QuoteCurrency, "tst")
	}
}
