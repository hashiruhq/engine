package engine_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/around25/products/matching-engine/model"
)

func TestMarketConvertFromAndToBinary(t *testing.T) {
	Convey("Should be able to convert a market from/to binary data", t, func() {
		market := model.MarketBackup{Topic: "orders", Partition: 1, Offset: 1213, HighestBid: 12200000000, LowestAsk: 12231000000}
		data, _ := market.ToBinary()
		market.FromBinary(data)
		So(market.GetTopic(), ShouldEqual, "orders")
		So(market.GetPartition(), ShouldEqual, 1)
		So(market.GetOffset(), ShouldEqual, 1213)
		bytes, _ := market.ToBinary()
		So(bytes, ShouldResemble, data)
	})
}
