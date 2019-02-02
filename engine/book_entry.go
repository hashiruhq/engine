package engine

import proto "github.com/golang/protobuf/proto"

// FromBinary loads a book entry from a byte array
func (bookEntry *BookEntry) FromBinary(msg []byte) error {
	return proto.Unmarshal(msg, bookEntry)
}

// ToBinary converts a book entry to a byte string
func (bookEntry *BookEntry) ToBinary() ([]byte, error) {
	return proto.Marshal(bookEntry)
}

// BookEntries defines convert operations on book entries
type BookEntries []BookEntry
