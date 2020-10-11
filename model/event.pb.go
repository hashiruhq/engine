// Code generated by protoc-gen-go. DO NOT EDIT.
// source: event.proto

package model

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

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
	// Error in processing
	EventType_Error EventType = 4
)

var EventType_name = map[int32]string{
	0: "Unspecified",
	1: "OrderStatusChange",
	2: "NewTrade",
	3: "OrderActivated",
	4: "Error",
}

var EventType_value = map[string]int32{
	"Unspecified":       0,
	"OrderStatusChange": 1,
	"NewTrade":          2,
	"OrderActivated":    3,
	"Error":             4,
}

func (x EventType) String() string {
	return proto.EnumName(EventType_name, int32(x))
}

func (EventType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{0}
}

type ErrorCode int32

const (
	ErrorCode_Undefined    ErrorCode = 0
	ErrorCode_InvalidOrder ErrorCode = 1
	ErrorCode_CancelFailed ErrorCode = 2
)

var ErrorCode_name = map[int32]string{
	0: "Undefined",
	1: "InvalidOrder",
	2: "CancelFailed",
}

var ErrorCode_value = map[string]int32{
	"Undefined":    0,
	"InvalidOrder": 1,
	"CancelFailed": 2,
}

func (x ErrorCode) String() string {
	return proto.EnumName(ErrorCode_name, int32(x))
}

func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{1}
}

type OrderStatusMsg struct {
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

func (m *OrderStatusMsg) Reset()         { *m = OrderStatusMsg{} }
func (m *OrderStatusMsg) String() string { return proto.CompactTextString(m) }
func (*OrderStatusMsg) ProtoMessage()    {}
func (*OrderStatusMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{0}
}

func (m *OrderStatusMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderStatusMsg.Unmarshal(m, b)
}
func (m *OrderStatusMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderStatusMsg.Marshal(b, m, deterministic)
}
func (m *OrderStatusMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderStatusMsg.Merge(m, src)
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

func (m *OrderStatusMsg) GetType() OrderType {
	if m != nil {
		return m.Type
	}
	return OrderType_Limit
}

func (m *OrderStatusMsg) GetSide() MarketSide {
	if m != nil {
		return m.Side
	}
	return MarketSide_Buy
}

func (m *OrderStatusMsg) GetPrice() uint64 {
	if m != nil {
		return m.Price
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

func (m *OrderStatusMsg) GetOwnerID() uint64 {
	if m != nil {
		return m.OwnerID
	}
	return 0
}

func (m *OrderStatusMsg) GetStatus() OrderStatus {
	if m != nil {
		return m.Status
	}
	return OrderStatus_Pending
}

type ErrorMsg struct {
	Code ErrorCode `protobuf:"varint,1,opt,name=Code,proto3,enum=model.ErrorCode" json:"Code,omitempty"`
	// The type of the order: 0=limit 1=market
	Type OrderType `protobuf:"varint,2,opt,name=Type,proto3,enum=model.OrderType" json:"Type,omitempty"`
	// The side of the market: 0=buy 1=sell
	Side    MarketSide `protobuf:"varint,3,opt,name=Side,proto3,enum=model.MarketSide" json:"Side,omitempty"`
	OrderID uint64     `protobuf:"varint,4,opt,name=OrderID,proto3" json:"OrderID,omitempty"`
	Price   uint64     `protobuf:"varint,5,opt,name=Price,proto3" json:"Price,omitempty"`
	Amount  uint64     `protobuf:"varint,6,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Funds   uint64     `protobuf:"varint,7,opt,name=Funds,proto3" json:"Funds,omitempty"`
	// The unique identifier the account that added the order
	OwnerID              uint64   `protobuf:"varint,8,opt,name=OwnerID,proto3" json:"OwnerID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrorMsg) Reset()         { *m = ErrorMsg{} }
func (m *ErrorMsg) String() string { return proto.CompactTextString(m) }
func (*ErrorMsg) ProtoMessage()    {}
func (*ErrorMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{1}
}

func (m *ErrorMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrorMsg.Unmarshal(m, b)
}
func (m *ErrorMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrorMsg.Marshal(b, m, deterministic)
}
func (m *ErrorMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorMsg.Merge(m, src)
}
func (m *ErrorMsg) XXX_Size() int {
	return xxx_messageInfo_ErrorMsg.Size(m)
}
func (m *ErrorMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorMsg.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorMsg proto.InternalMessageInfo

func (m *ErrorMsg) GetCode() ErrorCode {
	if m != nil {
		return m.Code
	}
	return ErrorCode_Undefined
}

func (m *ErrorMsg) GetType() OrderType {
	if m != nil {
		return m.Type
	}
	return OrderType_Limit
}

func (m *ErrorMsg) GetSide() MarketSide {
	if m != nil {
		return m.Side
	}
	return MarketSide_Buy
}

func (m *ErrorMsg) GetOrderID() uint64 {
	if m != nil {
		return m.OrderID
	}
	return 0
}

func (m *ErrorMsg) GetPrice() uint64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *ErrorMsg) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *ErrorMsg) GetFunds() uint64 {
	if m != nil {
		return m.Funds
	}
	return 0
}

func (m *ErrorMsg) GetOwnerID() uint64 {
	if m != nil {
		return m.OwnerID
	}
	return 0
}

type Event struct {
	Type      EventType `protobuf:"varint,1,opt,name=Type,proto3,enum=model.EventType" json:"Type,omitempty"`
	Market    string    `protobuf:"bytes,2,opt,name=Market,proto3" json:"Market,omitempty"`
	CreatedAt int64     `protobuf:"varint,3,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	// Types that are valid to be assigned to Payload:
	//	*Event_OrderStatus
	//	*Event_Trade
	//	*Event_OrderActivation
	//	*Event_Error
	Payload              isEvent_Payload `protobuf_oneof:"Payload"`
	SeqID                uint64          `protobuf:"varint,7,opt,name=SeqID,proto3" json:"SeqID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{2}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
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
	OrderActivation *OrderStatusMsg `protobuf:"bytes,6,opt,name=OrderActivation,proto3,oneof"`
}

type Event_Error struct {
	Error *ErrorMsg `protobuf:"bytes,8,opt,name=Error,proto3,oneof"`
}

func (*Event_OrderStatus) isEvent_Payload() {}

func (*Event_Trade) isEvent_Payload() {}

func (*Event_OrderActivation) isEvent_Payload() {}

func (*Event_Error) isEvent_Payload() {}

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

func (m *Event) GetOrderActivation() *OrderStatusMsg {
	if x, ok := m.GetPayload().(*Event_OrderActivation); ok {
		return x.OrderActivation
	}
	return nil
}

func (m *Event) GetError() *ErrorMsg {
	if x, ok := m.GetPayload().(*Event_Error); ok {
		return x.Error
	}
	return nil
}

func (m *Event) GetSeqID() uint64 {
	if m != nil {
		return m.SeqID
	}
	return 0
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Event) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Event_OrderStatus)(nil),
		(*Event_Trade)(nil),
		(*Event_OrderActivation)(nil),
		(*Event_Error)(nil),
	}
}

type Events struct {
	Events               []*Event `protobuf:"bytes,1,rep,name=Events,proto3" json:"Events,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Events) Reset()         { *m = Events{} }
func (m *Events) String() string { return proto.CompactTextString(m) }
func (*Events) ProtoMessage()    {}
func (*Events) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{3}
}

func (m *Events) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Events.Unmarshal(m, b)
}
func (m *Events) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Events.Marshal(b, m, deterministic)
}
func (m *Events) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Events.Merge(m, src)
}
func (m *Events) XXX_Size() int {
	return xxx_messageInfo_Events.Size(m)
}
func (m *Events) XXX_DiscardUnknown() {
	xxx_messageInfo_Events.DiscardUnknown(m)
}

var xxx_messageInfo_Events proto.InternalMessageInfo

func (m *Events) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

func init() {
	proto.RegisterEnum("model.EventType", EventType_name, EventType_value)
	proto.RegisterEnum("model.ErrorCode", ErrorCode_name, ErrorCode_value)
	proto.RegisterType((*OrderStatusMsg)(nil), "model.OrderStatusMsg")
	proto.RegisterType((*ErrorMsg)(nil), "model.ErrorMsg")
	proto.RegisterType((*Event)(nil), "model.Event")
	proto.RegisterType((*Events)(nil), "model.Events")
}

func init() { proto.RegisterFile("event.proto", fileDescriptor_2d17a9d3f0ddf27e) }

var fileDescriptor_2d17a9d3f0ddf27e = []byte{
	// 528 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xc1, 0x6e, 0xda, 0x40,
	0x10, 0x8d, 0x8d, 0x6d, 0xf0, 0x38, 0x0d, 0xce, 0xaa, 0x89, 0xac, 0xa8, 0x07, 0x84, 0xa8, 0x1a,
	0x71, 0xe0, 0x40, 0x4f, 0xbd, 0x54, 0xa2, 0x90, 0x28, 0x1c, 0x68, 0x22, 0x93, 0xdc, 0xbb, 0x65,
	0x27, 0xd4, 0x2a, 0x59, 0xd3, 0x65, 0x21, 0xca, 0x67, 0xf6, 0x63, 0xaa, 0x5e, 0xab, 0x9d, 0xdd,
	0x80, 0x9b, 0xa6, 0x3d, 0xf5, 0xc6, 0x7b, 0x6f, 0x66, 0x98, 0xf7, 0x76, 0x00, 0x12, 0xdc, 0xa0,
	0xd4, 0xbd, 0xa5, 0x2a, 0x75, 0xc9, 0xc2, 0xbb, 0x52, 0xe0, 0xe2, 0x24, 0x29, 0x95, 0x40, 0x65,
	0xb9, 0x93, 0x44, 0x2b, 0x2e, 0xd0, 0x82, 0xf6, 0x4f, 0x0f, 0x0e, 0x2e, 0x8d, 0x38, 0xd5, 0x5c,
	0xaf, 0x57, 0x93, 0xd5, 0x9c, 0x1d, 0x80, 0x3f, 0x1e, 0x65, 0x5e, 0xcb, 0x3b, 0x0d, 0x72, 0x7f,
	0x3c, 0x62, 0x1d, 0x08, 0xae, 0x1f, 0x96, 0x98, 0xf9, 0x2d, 0xef, 0xf4, 0xa0, 0x9f, 0xf6, 0x68,
	0x64, 0x8f, 0x9a, 0x0c, 0x9f, 0x93, 0xca, 0x5e, 0x43, 0x30, 0x2d, 0x04, 0x66, 0x35, 0xaa, 0x3a,
	0x74, 0x55, 0x13, 0xae, 0xbe, 0xa2, 0x36, 0x42, 0x4e, 0x32, 0x7b, 0x09, 0xe1, 0x95, 0x2a, 0x66,
	0x98, 0x05, 0x34, 0xdf, 0x02, 0x76, 0x0c, 0xd1, 0xe0, 0xae, 0x5c, 0x4b, 0x9d, 0x85, 0x44, 0x3b,
	0x64, 0xaa, 0xcf, 0xd7, 0x52, 0xac, 0xb2, 0xc8, 0x56, 0x13, 0x60, 0x19, 0xd4, 0x2f, 0xef, 0x25,
	0xaa, 0xf1, 0x28, 0xab, 0x13, 0xff, 0x08, 0x59, 0x17, 0x22, 0xeb, 0x23, 0x6b, 0xd0, 0x1a, 0xac,
	0xba, 0xac, 0x55, 0x72, 0x57, 0xd1, 0xfe, 0xe1, 0x41, 0xe3, 0x4c, 0xa9, 0x52, 0x19, 0xcf, 0x1d,
	0x08, 0x86, 0xa5, 0x40, 0x72, 0xbd, 0xf3, 0x48, 0xb2, 0xe1, 0x73, 0x52, 0xff, 0x6f, 0x12, 0xc6,
	0x85, 0xe9, 0x1c, 0x8f, 0x5c, 0x16, 0x8f, 0x70, 0x97, 0x51, 0xf8, 0x7c, 0x46, 0xd1, 0xf3, 0x19,
	0xd5, 0xff, 0x92, 0x51, 0xe3, 0xb7, 0x8c, 0xda, 0xdf, 0x7d, 0x08, 0xcf, 0xcc, 0x89, 0x6c, 0xed,
	0x3c, 0x31, 0x6d, 0xb4, 0x8a, 0x9d, 0x63, 0x88, 0xec, 0xee, 0x64, 0x3b, 0xce, 0x1d, 0x62, 0xaf,
	0x20, 0x1e, 0x2a, 0xe4, 0x1a, 0xc5, 0x40, 0x93, 0xd7, 0x5a, 0xbe, 0x23, 0xd8, 0x3b, 0x48, 0x2a,
	0xa1, 0x93, 0xc3, 0xa4, 0x7f, 0xf4, 0xe7, 0x73, 0x4c, 0x56, 0xf3, 0x8b, 0xbd, 0xbc, 0x5a, 0xcb,
	0x3a, 0x10, 0x5e, 0x9b, 0x0b, 0x25, 0xfb, 0x49, 0x7f, 0xdf, 0x35, 0x11, 0x77, 0xb1, 0x97, 0x5b,
	0x91, 0x0d, 0xa0, 0x49, 0x4d, 0x83, 0x99, 0x2e, 0x36, 0x5c, 0x17, 0xa5, 0xa4, 0x5c, 0xfe, 0xf1,
	0x25, 0x4f, 0xeb, 0xd9, 0x1b, 0x08, 0xe9, 0x85, 0x29, 0xa1, 0xa4, 0xdf, 0xac, 0xbe, 0xba, 0x6d,
	0xb1, 0xba, 0x89, 0x78, 0x8a, 0xdf, 0xb6, 0xe7, 0x66, 0xc1, 0x87, 0x18, 0xea, 0x57, 0xfc, 0x61,
	0x51, 0x72, 0xd1, 0xee, 0x41, 0x44, 0xb1, 0x99, 0xe5, 0xdd, 0xa7, 0xcc, 0x6b, 0xd5, 0x2a, 0xdb,
	0x13, 0x99, 0x3b, 0xad, 0xfb, 0x09, 0xe2, 0x6d, 0xcc, 0xac, 0x09, 0xc9, 0x8d, 0x5c, 0x2d, 0x71,
	0x56, 0xdc, 0x16, 0x28, 0xd2, 0x3d, 0x76, 0x04, 0x87, 0x95, 0xe5, 0x87, 0x5f, 0xb8, 0x9c, 0x63,
	0xea, 0xb1, 0x7d, 0x68, 0x7c, 0xc4, 0x7b, 0x72, 0x9f, 0xfa, 0x8c, 0xb9, 0xdf, 0xad, 0xf3, 0x83,
	0x22, 0xad, 0xb1, 0xd8, 0x19, 0x4a, 0x83, 0xee, 0x7b, 0x88, 0xb7, 0xd7, 0xcb, 0x5e, 0x40, 0x7c,
	0x23, 0x05, 0xde, 0x16, 0x92, 0xe6, 0xa7, 0xb0, 0x3f, 0x96, 0x1b, 0xbe, 0x28, 0x04, 0x4d, 0x48,
	0x3d, 0xc3, 0x0c, 0xb9, 0x9c, 0xe1, 0xe2, 0x9c, 0x17, 0x0b, 0x14, 0xa9, 0xff, 0x39, 0xa2, 0xbf,
	0x87, 0xb7, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x52, 0xb9, 0x93, 0xc6, 0x4e, 0x04, 0x00, 0x00,
}
