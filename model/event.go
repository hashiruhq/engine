package model

import (
	"time"

	proto "github.com/golang/protobuf/proto"
)

// NewOrderStatusEvent returns a new event set with the order status details
func NewOrderStatusEvent(market string, id uint64, amount, funds uint64, status OrderStatus) Event {
	return Event{
		Type:   EventType_OrderStatusChange,
		Market: market,
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

// NewOrderActivatedEvent returns a new event set with the order status details
func NewOrderActivatedEvent(market string, orderType OrderType, side MarketSide, id, ownerID uint64, price, amount, funds uint64, status OrderStatus) Event {
	return Event{
		Type:   EventType_OrderActivated,
		Market: market,
		Payload: &Event_OrderActivation{
			OrderActivation: &StopOrderActivatedMsg{
				ID:      id,
				Type:    orderType,
				Side:    side,
				OwnerID: ownerID,
				Amount:  amount,
				Funds:   funds,
				Price:   price,
				Status:  status,
			},
		},
		CreatedAt: time.Now().UTC().Unix(),
	}
}

// NewTradeEvent returns a new event set with the trade details
func NewTradeEvent(market string, takerSide MarketSide, askID, bidID, askOwnerID, bidOwnerID, amount, price uint64) Event {
	return Event{
		Type:   EventType_NewTrade,
		Market: market,
		Payload: &Event_Trade{
			Trade: &Trade{
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
