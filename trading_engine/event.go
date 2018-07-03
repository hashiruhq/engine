package trading_engine

import (
	"github.com/Shopify/sarama"
)

type Event struct {
	Msg    *sarama.ConsumerMessage
	Order  Order
	Trades []Trade
}

// NewEvent Create a new event
func NewEvent(msg *sarama.ConsumerMessage) Event {
	return Event{Msg: msg}
}

func (event *Event) Decode() {
	event.Order.FromJSON(event.Msg.Value)
}

func (event *Event) SetTrades(trades []Trade) {
	event.Trades = trades
}

func (event Event) HasTrades() bool {
	return len(event.Trades) > 0
}
