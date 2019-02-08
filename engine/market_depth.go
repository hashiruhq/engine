package engine

import proto "github.com/golang/protobuf/proto"

// FromBinary loads a market backup from a byte array
func (m *MarketDepth) FromBinary(msg []byte) error {
	return proto.Unmarshal(msg, m)
}

// ToBinary converts a market backup to a byte string
func (m *MarketDepth) ToBinary() ([]byte, error) {
	return proto.Marshal(m)
}
