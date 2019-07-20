package model

import proto "github.com/golang/protobuf/proto"

// FromBinary loads a market backup from a byte array
func (m *MarketBackup) FromBinary(msg []byte) error {
	return proto.Unmarshal(msg, m)
}

// ToBinary converts a market backup to a byte string
func (m *MarketBackup) ToBinary() ([]byte, error) {
	return proto.Marshal(m)
}
