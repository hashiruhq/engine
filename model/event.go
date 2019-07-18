package model

import (
	"time"

	proto "github.com/golang/protobuf/proto"
)

// NewOrderStatusEvent returns a new event set with the order status details
func NewOrderStatusEvent(id uint64, amount, funds uint64, status OrderStatus) Event {
	return Event{
		Type: EventType_OrderStatusChange,
		Payload: &Event_OrderStatus{
			OrderStatus: &OrderStatusMsg{
				ID:     id,
				Amount: amount,
				Funds:  funds,
				Status: status,
			},
		},
		CreatedAt: time.Now().UTC().Unix(),
	}
}

// NewTradeEvent returns a new event set with the trade details
func NewTradeEvent(market string, takerSide MarketSide, askID, bidID, askOwnerID, bidOwnerID, amount, price uint64) Event {
	return Event{
		Type: EventType_NewTrade,
		Payload: &Event_Trade{
			Trade: &Trade{
				Market:     market,
				TakerSide:  takerSide,
				AskID:      askID,
				AskOwnerID: askOwnerID,
				BidID:      bidID,
				BidOwnerID: bidOwnerID,
				Amount:     amount,
				Price:      price,
			},
		},
		CreatedAt: time.Now().UTC().Unix(),
	}
}

// FromBinary loads an event from a byte array
func (event *Event) FromBinary(msg []byte) error {
	return proto.Unmarshal(msg, event)
}

// ToBinary converts an event to a byte string
func (event *Event) ToBinary() ([]byte, error) {
	return proto.Marshal(event)
}
