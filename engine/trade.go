package engine

import proto "github.com/golang/protobuf/proto"

// NewTrade Creates a new trade between the taker order and the maker order
func NewTrade(takerOrder string, makerOrder string, amount, price uint64) Trade {
	return Trade{
		TakerOrderID: takerOrder,
		MakerOrderID: makerOrder,
		Amount:       amount,
		Price:        price,
	}
}

// FromBinary loads a trade from a byte array
func (trade *Trade) FromBinary(msg []byte) error {
	return proto.Unmarshal(msg, trade)
}

// ToBinary converts a trade to a byte string
func (trade *Trade) ToBinary() ([]byte, error) {
	return proto.Marshal(trade)
}
