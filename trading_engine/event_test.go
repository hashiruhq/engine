package trading_engine_test

import (
	"testing"

	"gitlab.com/around25/products/matching-engine/trading_engine"

	"github.com/Shopify/sarama"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEventsUsage(t *testing.T) {
	Convey("Create a new event", t, func() {
		msg := &sarama.ConsumerMessage{Value: []byte(`{"id": "112", "price": "12.00", "amount": "1213.00"}`)}
		event := trading_engine.NewEvent(msg)
		Convey("I should be able to decode the message as an order", func() {
			event.Decode()
			So(event.Order.ID, ShouldEqual, "112")
			So(event.Order.Price, ShouldEqual, 1200000000)
			So(event.Order.Amount, ShouldEqual, 121300000000)
			So(event.HasTrades(), ShouldEqual, false)
			trades := make([]trading_engine.Trade, 2)
			event.SetTrades(trades)
			So(event.HasTrades(), ShouldEqual, true)
		})
	})
}
