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

func TestTradeLoadFromBinary(t *testing.T) {
	Convey("Should be able to load a trade from binary", t, func() {
		trade := NewTrade("TST_1", "TST_2", 848382829993942, 131221300010201)
		data, _ := trade.ToBinary()
		trade = Trade{}
		trade.FromBinary(data)
		So(trade.Price, ShouldEqual, 131221300010201)
		So(trade.Amount, ShouldEqual, 848382829993942)
		So(trade.TakerOrderID, ShouldEqual, "TST_1")
		So(trade.MakerOrderID, ShouldEqual, "TST_2")
	})
}

func TestTradeConvertToBinary(t *testing.T) {
	Convey("Should be able to convert a trade to binary", t, func() {
		trade := NewTrade("TST_1", "TST_2", 848382829993942, 131221300010201)
		initData, _ := trade.ToBinary()
		trade = Trade{}
		trade.FromBinary(initData)
		encoded, _ := trade.ToBinary()
		So(encoded, ShouldResemble, initData)
	})
}
