package trading_engine_test

import (
	"testing"

	"gitlab.com/around25/products/matching-engine/trading_engine"

	"github.com/francoispqt/gojay"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMarketConvertFrom_ToEmptyJson(t *testing.T) {
	Convey("Should be able to convert a market from/to empty json string", t, func() {
		var market trading_engine.Market
		json := ``
		market.FromJSON([]byte(json))
		bytes, _ := gojay.Marshal(market)
		So(string(bytes), ShouldEqual, json)
		So(market.IsNil(), ShouldEqual, true)
	})

	Convey("Should be able to convert a market from/to json string", t, func() {
		var market trading_engine.Market
		json := `{"topic":"orders","partition":1,"offset":1213,"lowest_ask":"122.31000000","highest_bid":"122.00000000","buy_orders":[],"sell_orders":[]}`
		market.FromJSON([]byte(json))
		So(market.Topic, ShouldEqual, "orders")
		So(market.Partition, ShouldEqual, 1)
		So(market.Offset, ShouldEqual, 1213)
		So(market.IsNil(), ShouldEqual, false)

		bytes, _ := market.ToJSON()
		So(string(bytes), ShouldEqual, json)
	})
}
