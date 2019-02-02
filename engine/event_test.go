package engine_test

import (
	"testing"

	"github.com/Shopify/sarama"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/around25/products/matching-engine/engine"
)

func TestEventsUsage(t *testing.T) {
	Convey("Create a new event", t, func() {
		order := engine.Order{ID: "112", Price: 1200000000, Amount: 121300000000}
		encoded, _ := order.ToBinary()
		msg := &sarama.ConsumerMessage{Value: encoded}
		event := engine.NewEvent(msg)
		Convey("I should be able to decode the message as an order", func() {
			event.Decode()
			So(event.Order.ID, ShouldEqual, "112")
			So(event.Order.Price, ShouldEqual, 1200000000)
			So(event.Order.Amount, ShouldEqual, 121300000000)
			So(event.HasTrades(), ShouldEqual, false)
			trades := make([]engine.Trade, 2)
			event.SetTrades(trades)
			So(event.HasTrades(), ShouldEqual, true)
		})
	})
}
