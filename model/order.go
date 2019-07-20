package model

import (
	proto "github.com/golang/protobuf/proto"
)

// NewOrder create a new order
func NewOrder(id, price, amount uint64, side MarketSide, category OrderType, eventType CommandType) Order {
	return Order{ID: id, Price: price, Amount: amount, Side: side, Type: category, EventType: eventType}
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
