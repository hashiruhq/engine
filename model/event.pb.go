// Code generated by protoc-gen-go. DO NOT EDIT.
// source: event.proto

package model

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EventType int32

const (
	// Invalid status for any non set field
	EventType_Unspecified EventType = 0
	// An order has changed status, either by being acknowledged by the engine on creation,
	// being cancelled by a command or by matching with another order.
	EventType_OrderStatusChange EventType = 1
	// A new trade was generated beetween two orders
	EventType_NewTrade EventType = 2
	// A stop loss order was activated by the engine due to market change
	EventType_OrderActivated EventType = 3
)

var EventType_name = map[int32]string{
	0: "Unspecified",
	1: "OrderStatusChange",
	2: "NewTrade",
	3: "OrderActivated",
}
var EventType_value = map[string]int32{
	"Unspecified":       0,
	"OrderStatusChange": 1,
	"NewTrade":          2,
	"OrderActivated":    3,
}

func (x EventType) String() string {
	return proto.EnumName(EventType_name, int32(x))
}
func (EventType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_event_e488d22eb2081aa4, []int{0}
}

type OrderStatusMsg struct {
	ID                   uint64      `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Amount               uint64      `protobuf:"varint,2,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Funds                uint64      `protobuf:"varint,3,opt,name=Funds,proto3" json:"Funds,omitempty"`
	Status               OrderStatus `protobuf:"varint,4,opt,name=Status,proto3,enum=model.OrderStatus" json:"Status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *OrderStatusMsg) Reset()         { *m = OrderStatusMsg{} }
func (m *OrderStatusMsg) String() string { return proto.CompactTextString(m) }
func (*OrderStatusMsg) ProtoMessage()    {}
func (*OrderStatusMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_event_e488d22eb2081aa4, []int{0}
}
func (m *OrderStatusMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderStatusMsg.Unmarshal(m, b)
}
func (m *OrderStatusMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderStatusMsg.Marshal(b, m, deterministic)
}
func (dst *OrderStatusMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderStatusMsg.Merge(dst, src)
}
func (m *OrderStatusMsg) XXX_Size() int {
	return xxx_messageInfo_OrderStatusMsg.Size(m)
}
func (m *OrderStatusMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderStatusMsg.DiscardUnknown(m)
}

var xxx_messageInfo_OrderStatusMsg proto.InternalMessageInfo

func (m *OrderStatusMsg) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *OrderStatusMsg) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *OrderStatusMsg) GetFunds() uint64 {
	if m != nil {
		return m.Funds
	}
	return 0
}

func (m *OrderStatusMsg) GetStatus() OrderStatus {
	if m != nil {
		return m.Status
	}
	return OrderStatus_Pending
}

type StopOrderActivatedMsg struct {
	ID uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	// The type of the order: 0=limit 1=market
	Type OrderType `protobuf:"varint,2,opt,name=Type,proto3,enum=model.OrderType" json:"Type,omitempty"`
	// The side of the market: 0=buy 1=sell
	Side   MarketSide `protobuf:"varint,3,opt,name=Side,proto3,enum=model.MarketSide" json:"Side,omitempty"`
	Price  uint64     `protobuf:"varint,4,opt,name=Price,proto3" json:"Price,omitempty"`
	Amount uint64     `protobuf:"varint,5,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Funds  uint64     `protobuf:"varint,6,opt,name=Funds,proto3" json:"Funds,omitempty"`
	// The unique identifier the account that added the order
	OwnerID              uint64      `protobuf:"varint,7,opt,name=OwnerID,proto3" json:"OwnerID,omitempty"`
	Status               OrderStatus `protobuf:"varint,8,opt,name=Status,proto3,enum=model.OrderStatus" json:"Status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *StopOrderActivatedMsg) Reset()         { *m = StopOrderActivatedMsg{} }
func (m *StopOrderActivatedMsg) String() string { return proto.CompactTextString(m) }
func (*StopOrderActivatedMsg) ProtoMessage()    {}
func (*StopOrderActivatedMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_event_e488d22eb2081aa4, []int{1}
}
func (m *StopOrderActivatedMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopOrderActivatedMsg.Unmarshal(m, b)
}
func (m *StopOrderActivatedMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopOrderActivatedMsg.Marshal(b, m, deterministic)
}
func (dst *StopOrderActivatedMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopOrderActivatedMsg.Merge(dst, src)
}
func (m *StopOrderActivatedMsg) XXX_Size() int {
	return xxx_messageInfo_StopOrderActivatedMsg.Size(m)
}
func (m *StopOrderActivatedMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_StopOrderActivatedMsg.DiscardUnknown(m)
}

var xxx_messageInfo_StopOrderActivatedMsg proto.InternalMessageInfo

func (m *StopOrderActivatedMsg) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *StopOrderActivatedMsg) GetType() OrderType {
	if m != nil {
		return m.Type
	}
	return OrderType_Limit
}

func (m *StopOrderActivatedMsg) GetSide() MarketSide {
	if m != nil {
		return m.Side
	}
	return MarketSide_Buy
}

func (m *StopOrderActivatedMsg) GetPrice() uint64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *StopOrderActivatedMsg) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *StopOrderActivatedMsg) GetFunds() uint64 {
	if m != nil {
		return m.Funds
	}
	return 0
}

func (m *StopOrderActivatedMsg) GetOwnerID() uint64 {
	if m != nil {
		return m.OwnerID
	}
	return 0
}

func (m *StopOrderActivatedMsg) GetStatus() OrderStatus {
	if m != nil {
		return m.Status
	}
	return OrderStatus_Pending
}

type Event struct {
	Type      EventType `protobuf:"varint,1,opt,name=Type,proto3,enum=model.EventType" json:"Type,omitempty"`
	Market    string    `protobuf:"bytes,2,opt,name=Market,proto3" json:"Market,omitempty"`
	CreatedAt int64     `protobuf:"varint,3,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	// Types that are valid to be assigned to Payload:
	//	*Event_OrderStatus
	//	*Event_Trade
	//	*Event_OrderActivation
	Payload              isEvent_Payload `protobuf_oneof:"Payload"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_event_e488d22eb2081aa4, []int{2}
}
func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (dst *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(dst, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetType() EventType {
	if m != nil {
		return m.Type
	}
	return EventType_Unspecified
}

func (m *Event) GetMarket() string {
	if m != nil {
		return m.Market
	}
	return ""
}

func (m *Event) GetCreatedAt() int64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

type isEvent_Payload interface {
	isEvent_Payload()
}

type Event_OrderStatus struct {
	OrderStatus *OrderStatusMsg `protobuf:"bytes,4,opt,name=OrderStatus,proto3,oneof"`
}

type Event_Trade struct {
	Trade *Trade `protobuf:"bytes,5,opt,name=Trade,proto3,oneof"`
}

type Event_OrderActivation struct {
	OrderActivation *StopOrderActivatedMsg `protobuf:"bytes,6,opt,name=OrderActivation,proto3,oneof"`
}

func (*Event_OrderStatus) isEvent_Payload() {}

func (*Event_Trade) isEvent_Payload() {}

func (*Event_OrderActivation) isEvent_Payload() {}

func (m *Event) GetPayload() isEvent_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Event) GetOrderStatus() *OrderStatusMsg {
	if x, ok := m.GetPayload().(*Event_OrderStatus); ok {
		return x.OrderStatus
	}
	return nil
}

func (m *Event) GetTrade() *Trade {
	if x, ok := m.GetPayload().(*Event_Trade); ok {
		return x.Trade
	}
	return nil
}

func (m *Event) GetOrderActivation() *StopOrderActivatedMsg {
	if x, ok := m.GetPayload().(*Event_OrderActivation); ok {
		return x.OrderActivation
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Event) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Event_OneofMarshaler, _Event_OneofUnmarshaler, _Event_OneofSizer, []interface{}{
		(*Event_OrderStatus)(nil),
		(*Event_Trade)(nil),
		(*Event_OrderActivation)(nil),
	}
}

func _Event_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Event)
	// Payload
	switch x := m.Payload.(type) {
	case *Event_OrderStatus:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.OrderStatus); err != nil {
			return err
		}
	case *Event_Trade:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Trade); err != nil {
			return err
		}
	case *Event_OrderActivation:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.OrderActivation); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Event.Payload has unexpected type %T", x)
	}
	return nil
}

func _Event_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Event)
	switch tag {
	case 4: // Payload.OrderStatus
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(OrderStatusMsg)
		err := b.DecodeMessage(msg)
		m.Payload = &Event_OrderStatus{msg}
		return true, err
	case 5: // Payload.Trade
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Trade)
		err := b.DecodeMessage(msg)
		m.Payload = &Event_Trade{msg}
		return true, err
	case 6: // Payload.OrderActivation
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(StopOrderActivatedMsg)
		err := b.DecodeMessage(msg)
		m.Payload = &Event_OrderActivation{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Event_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Event)
	// Payload
	switch x := m.Payload.(type) {
	case *Event_OrderStatus:
		s := proto.Size(x.OrderStatus)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Trade:
		s := proto.Size(x.Trade)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_OrderActivation:
		s := proto.Size(x.OrderActivation)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*OrderStatusMsg)(nil), "model.OrderStatusMsg")
	proto.RegisterType((*StopOrderActivatedMsg)(nil), "model.StopOrderActivatedMsg")
	proto.RegisterType((*Event)(nil), "model.Event")
	proto.RegisterEnum("model.EventType", EventType_name, EventType_value)
}

func init() { proto.RegisterFile("event.proto", fileDescriptor_event_e488d22eb2081aa4) }

var fileDescriptor_event_e488d22eb2081aa4 = []byte{
	// 419 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x8d, 0x1d, 0xdb, 0x69, 0x66, 0x2b, 0xd7, 0x1d, 0x91, 0xca, 0xaa, 0x7a, 0xa8, 0xaa, 0x22,
	0x55, 0x3d, 0xe4, 0x60, 0x4e, 0x1c, 0x43, 0x03, 0x4a, 0x0e, 0xa1, 0xd5, 0xa6, 0xfd, 0x00, 0x93,
	0x1d, 0x82, 0x45, 0xeb, 0xb5, 0xd6, 0x9b, 0x56, 0xe1, 0x4f, 0xf8, 0x59, 0x84, 0x3c, 0x36, 0x61,
	0x81, 0x48, 0x1c, 0xdf, 0x7b, 0xb3, 0xbb, 0xef, 0xbd, 0x1d, 0x10, 0xf4, 0x4c, 0xa5, 0x1d, 0x57,
	0x46, 0x5b, 0x8d, 0xe1, 0x93, 0x56, 0xf4, 0x78, 0x2a, 0xb4, 0x51, 0x64, 0x5a, 0xee, 0x54, 0x58,
	0x93, 0x2b, 0x6a, 0xc1, 0xc5, 0x37, 0x88, 0x6f, 0x1b, 0x6d, 0x69, 0x73, 0xbb, 0xa9, 0x17, 0xf5,
	0x1a, 0x63, 0xf0, 0xe7, 0xd3, 0xd4, 0x3b, 0xf7, 0xae, 0x02, 0xe9, 0xcf, 0xa7, 0x78, 0x02, 0xd1,
	0xe4, 0x49, 0x6f, 0x4a, 0x9b, 0xfa, 0xcc, 0x75, 0x08, 0x5f, 0x41, 0xf8, 0x61, 0x53, 0xaa, 0x3a,
	0xed, 0x33, 0xdd, 0x02, 0xbc, 0x86, 0xa8, 0xbd, 0x2a, 0x0d, 0xce, 0xbd, 0xab, 0x38, 0xc3, 0x31,
	0x3b, 0x18, 0x3b, 0x8f, 0xc8, 0x6e, 0xe2, 0xe2, 0x87, 0x07, 0xa3, 0xa5, 0xd5, 0x15, 0x6b, 0x93,
	0x95, 0x2d, 0x9e, 0x73, 0x4b, 0x6a, 0x9f, 0x87, 0x4b, 0x08, 0xee, 0xb7, 0x15, 0xb1, 0x83, 0x38,
	0x4b, 0xdc, 0x3b, 0x1b, 0x5e, 0xb2, 0x8a, 0xaf, 0x21, 0x58, 0x16, 0x8a, 0xd8, 0x50, 0x9c, 0x1d,
	0x77, 0x53, 0x8b, 0xdc, 0x7c, 0x25, 0xdb, 0x08, 0x92, 0xe5, 0xc6, 0xf8, 0x9d, 0x29, 0x56, 0xc4,
	0x0e, 0x03, 0xd9, 0x02, 0x27, 0x66, 0xb8, 0x3f, 0x66, 0xe4, 0xc6, 0x4c, 0x61, 0x70, 0xfb, 0x52,
	0x92, 0x99, 0x4f, 0xd3, 0x01, 0xf3, 0xbf, 0xa0, 0x53, 0xc0, 0xc1, 0x7f, 0x0b, 0xf8, 0xee, 0x43,
	0xf8, 0xbe, 0xf9, 0xad, 0x5d, 0x40, 0xef, 0x8f, 0x80, 0xac, 0x39, 0x01, 0x4f, 0x20, 0x6a, 0xd3,
	0x70, 0x11, 0x43, 0xd9, 0x21, 0x3c, 0x83, 0xe1, 0x8d, 0xa1, 0xa6, 0xbc, 0x89, 0xe5, 0xf4, 0x7d,
	0xf9, 0x9b, 0xc0, 0xb7, 0x20, 0x9c, 0xc7, 0x39, 0xb5, 0xc8, 0x46, 0xff, 0xda, 0x5a, 0xd4, 0xeb,
	0x59, 0x4f, 0xba, 0xb3, 0x78, 0x09, 0xe1, 0x7d, 0xb3, 0x2c, 0xdc, 0x89, 0xc8, 0x0e, 0xbb, 0x43,
	0xcc, 0xcd, 0x7a, 0xb2, 0x15, 0x71, 0x06, 0x47, 0xee, 0x17, 0x16, 0xba, 0xe4, 0xb2, 0x44, 0x76,
	0xd6, 0xcd, 0xef, 0xfd, 0xe4, 0x59, 0x4f, 0xfe, 0x7d, 0xec, 0xdd, 0x10, 0x06, 0x77, 0xf9, 0xf6,
	0x51, 0xe7, 0xea, 0xfa, 0x01, 0x86, 0xbb, 0xf8, 0x78, 0x04, 0xe2, 0xa1, 0xac, 0x2b, 0x5a, 0x15,
	0x9f, 0x0b, 0x52, 0x49, 0x0f, 0x47, 0x70, 0xec, 0xf8, 0xbc, 0xf9, 0x92, 0x97, 0x6b, 0x4a, 0x3c,
	0x3c, 0x84, 0x83, 0x8f, 0xf4, 0xc2, 0xae, 0x12, 0x1f, 0xb1, 0xdb, 0xed, 0xdd, 0xab, 0x49, 0xff,
	0x53, 0xc4, 0x6b, 0xff, 0xe6, 0x67, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf4, 0x5f, 0x07, 0x3e, 0x26,
	0x03, 0x00, 0x00,
}
