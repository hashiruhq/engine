package model

import (
	"time"

	proto "github.com/golang/protobuf/proto"
)

// NewOrderStatusEvent returns a new event set with the order status details
func NewOrderStatusEvent(seqID uint64, market string, orderType OrderType, side MarketSide, id, ownerID uint64, price, amount, funds uint64, status OrderStatus, filledAmount uint64, usedFunds uint64) Event {
	return Event{
		SeqID:  seqID,
		Type:   EventType_OrderStatusChange,
		Market: market,
		Payload: &Event_OrderStatus{
			OrderStatus: &OrderStatusMsg{
				ID:           id,
				Type:         orderType,
				Side:         side,
				OwnerID:      ownerID,
				Price:        price,
				Amount:       amount,
				Funds:        funds,
				Status:       status,
				FilledAmount: filledAmount,
				UsedFunds:    usedFunds,
			},
		},
		CreatedAt: time.Now().UTC().UnixNano(),
	}
}

// NewOrderActivatedEvent returns a new event set with the order status details
func NewOrderActivatedEvent(seqID uint64, market string, orderType OrderType, side MarketSide, id, ownerID uint64, price, amount, funds uint64, status OrderStatus) Event {
	return Event{
		SeqID:  seqID,
		Type:   EventType_OrderActivated,
		Market: market,
		Payload: &Event_OrderActivation{
			OrderActivation: &OrderStatusMsg{
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
		CreatedAt: time.Now().UTC().UnixNano(),
	}
}

// NewTradeEvent returns a new event set with the trade details
func NewTradeEvent(seqID uint64, market string, tradeSeqID uint64, takerSide MarketSide, askID, bidID, askOwnerID, bidOwnerID, amount, price uint64) Event {
	return Event{
		SeqID:  seqID,
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
				SeqID:      tradeSeqID,
			},
		},
		CreatedAt: time.Now().UTC().UnixNano(),
	}
}

// NewErrorEvent returns a new error event
func NewErrorEvent(seqID uint64, market string, code ErrorCode, orderType OrderType, side MarketSide, id, ownerID, price, amount, funds uint64) Event {
	return Event{
		SeqID:  seqID,
		Type:   EventType_Error,
		Market: market,
		Payload: &Event_Error{
			Error: &ErrorMsg{
				Code:    code,
				OrderID: id,
				Type:    orderType,
				Side:    side,
				OwnerID: ownerID,
				Amount:  amount,
				Funds:   funds,
				Price:   price,
			},
		},
		CreatedAt: time.Now().UTC().UnixNano(),
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
