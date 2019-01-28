package engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTradeCreation(t *testing.T) {
	Convey("Should be able to create a new trade", t, func() {
		trade := NewTrade("ID_1", "ID_2", 12000000000, 100000000)
		So(trade.Amount, ShouldEqual, 12000000000)
		So(trade.Price, ShouldEqual, 100000000)
		So(trade.TakerOrderID, ShouldEqual, "ID_1")
		So(trade.MakerOrderID, ShouldEqual, "ID_2")
	})
}

func TestTradeLoadFromJson(t *testing.T) {
	Convey("Should be able to load a trade from json", t, func() {
		var trade Trade
		json := []byte(`{
			"taker_order_id":"TST_1",
			"maker_order_id":"TST_2",
			"price": "1312213.00010201",
			"amount": "8483828.29993942"
		}`)
		trade.FromJSON(json)
		So(trade.Price, ShouldEqual, 131221300010201)
		So(trade.Amount, ShouldEqual, 848382829993942)
		So(trade.TakerOrderID, ShouldEqual, "TST_1")
		So(trade.MakerOrderID, ShouldEqual, "TST_2")
	})
}

func TestTradeConvertToJson(t *testing.T) {
	Convey("Should be able to convert a trade to json string", t, func() {
		var trade Trade
		json := `{"taker_order_id":"TST_1","maker_order_id":"TST_2","price":"1312213.00010201","amount":"8483828.29993942"}`
		trade.FromJSON([]byte(json))
		bytes, _ := trade.ToJSON()
		So(string(bytes), ShouldEqual, json)
	})
}
