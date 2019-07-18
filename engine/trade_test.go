package engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTradeCreation(t *testing.T) {
	Convey("Should be able to create a new trade", t, func() {
		trade := NewTrade("btcusd", model.MarketSide_Buy, 1, 2, 1, 2, 12000000000, 100000000)
		So(trade.Amount, ShouldEqual, 12000000000)
		So(trade.Price, ShouldEqual, 100000000)
		So(trade.AskID, ShouldEqual, 1)
		So(trade.BidID, ShouldEqual, 2)
		So(trade.Market, ShouldEqual, "btcusd")
		So(trade.TakerSide, ShouldEqual, model.MarketSide_Buy)
	})
}

func TestTradeLoadFromBinary(t *testing.T) {
	Convey("Should be able to load a trade from binary", t, func() {
		trade := NewTrade("btcusd", model.MarketSide_Buy, 1, 2, 1, 2, 848382829993942, 131221300010201)
		data, _ := trade.ToBinary()
		trade = Trade{}
		trade.FromBinary(data)
		So(trade.Price, ShouldEqual, 131221300010201)
		So(trade.Amount, ShouldEqual, 848382829993942)
		So(trade.AskID, ShouldEqual, 1)
		So(trade.BidID, ShouldEqual, 2)
		So(trade.Market, ShouldEqual, "btcusd")
		So(trade.TakerSide, ShouldEqual, model.MarketSide_Buy)
	})
}

func TestTradeConvertToBinary(t *testing.T) {
	Convey("Should be able to convert a trade to binary", t, func() {
		trade := NewTrade("btcusd", model.MarketSide_Buy, 1, 2, 1, 2, 848382829993942, 131221300010201)
		initData, _ := trade.ToBinary()
		trade = Trade{}
		trade.FromBinary(initData)
		encoded, _ := trade.ToBinary()
		So(encoded, ShouldResemble, initData)
	})
}
