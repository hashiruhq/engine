package engine_test

import (
	"testing"

	"github.com/segmentio/kafka-go"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/around25/products/matching-engine/engine"
	"gitlab.com/around25/products/matching-engine/model"
)

func TestEventsUsage(t *testing.T) {
	Convey("Create a new event", t, func() {
		order := model.Order{ID: 1, Price: 1200000000, Amount: 121300000000}
		encoded, _ := order.ToBinary()
		msg := kafka.Message{Value: encoded}
		event := engine.NewEvent(msg)
		Convey("I should be able to decode the message as an order", func() {
			event.Decode()
			So(event.Order.ID, ShouldEqual, 1)
			So(event.Order.Price, ShouldEqual, 1200000000)
			So(event.Order.Amount, ShouldEqual, 121300000000)
			So(event.HasTrades(), ShouldEqual, false)
			trades := make([]model.Trade, 2)
			event.SetTrades(trades)
			So(event.HasTrades(), ShouldEqual, true)
		})
	})
}
