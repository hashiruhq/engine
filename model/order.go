package model

import (
	proto "github.com/golang/protobuf/proto"
)

// NewOrder create a new order
func NewOrder(id, price, amount uint64, side MarketSide, category OrderType, eventType CommandType) Order {
	return Order{ID: id, Price: price, Amount: amount, Side: side, Type: category, EventType: eventType}
}

// Valid checks if the order is valid based on the type of the order and the price/amount/funds
func (order *Order) Valid() bool {
	if order.ID == 0 {
		return false
	}
	switch order.EventType {
	case CommandType_NewOrder:
		{
			if order.Stop != StopLoss_None {
				if order.StopPrice == 0 {
					return false
				}
			}
			switch order.Type {
			case OrderType_Limit:
				return order.Price != 0 && order.Amount != 0
			case OrderType_Market:
				return order.Funds != 0 && order.Amount != 0
			}
		}
	case CommandType_CancelOrder:
		{
			if order.Stop != StopLoss_None {
				if order.StopPrice == 0 {
					return false
				}
			}
			if order.Type == OrderType_Limit {
				return order.Price != 0
			}
		}
	}
	return true
}

// Filled checks if the order can be considered filled
func (order *Order) Filled() bool {
	if order.EventType != CommandType_NewOrder {
		return false
	}
	if order.Type == OrderType_Market && (order.Amount == 0 || order.Funds == 0) {
		return true
	}
	return false
}

//***************************
// Interface Implementations
//***************************

// LessThan implementes the skiplist interface
func (order Order) LessThan(other Order) bool {
	return order.Price < other.Price
}

// FromBinary loads an order from a byte array
func (order *Order) FromBinary(msg []byte) error {
	return proto.Unmarshal(msg, order)
}

// ToBinary converts an order to a byte string
func (order *Order) ToBinary() ([]byte, error) {
	return proto.Marshal(order)
}

// SetStatus changes the status of an order if the new status has not already been set
// and it's not lower than the current set status
func (order *Order) SetStatus(status OrderStatus) {
	if order.Status < status {
		order.Status = status
	}
}
