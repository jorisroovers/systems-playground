// Code generated by protoc-gen-go. DO NOT EDIT.
// source: book_store.proto

package book_store

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Book struct {
	// Each field in the message definition has a:
	// 1) type
	// 2) name
	// 3) unique number: this is used during encoding, which allows for very efficient referencing to the fields
	Title                string   `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Author               string   `protobuf:"bytes,2,opt,name=author,proto3" json:"author,omitempty"`
	Year                 int32    `protobuf:"varint,3,opt,name=year,proto3" json:"year,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Book) Reset()         { *m = Book{} }
func (m *Book) String() string { return proto.CompactTextString(m) }
func (*Book) ProtoMessage()    {}
func (*Book) Descriptor() ([]byte, []int) {
	return fileDescriptor_book_store_1e7c586e95bda3a3, []int{0}
}
func (m *Book) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Book.Unmarshal(m, b)
}
func (m *Book) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Book.Marshal(b, m, deterministic)
}
func (dst *Book) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Book.Merge(dst, src)
}
func (m *Book) XXX_Size() int {
	return xxx_messageInfo_Book.Size(m)
}
func (m *Book) XXX_DiscardUnknown() {
	xxx_messageInfo_Book.DiscardUnknown(m)
}

var xxx_messageInfo_Book proto.InternalMessageInfo

func (m *Book) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Book) GetAuthor() string {
	if m != nil {
		return m.Author
	}
	return ""
}

func (m *Book) GetYear() int32 {
	if m != nil {
		return m.Year
	}
	return 0
}

type BookReference struct {
	Title                string   `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BookReference) Reset()         { *m = BookReference{} }
func (m *BookReference) String() string { return proto.CompactTextString(m) }
func (*BookReference) ProtoMessage()    {}
func (*BookReference) Descriptor() ([]byte, []int) {
	return fileDescriptor_book_store_1e7c586e95bda3a3, []int{1}
}
func (m *BookReference) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BookReference.Unmarshal(m, b)
}
func (m *BookReference) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BookReference.Marshal(b, m, deterministic)
}
func (dst *BookReference) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BookReference.Merge(dst, src)
}
func (m *BookReference) XXX_Size() int {
	return xxx_messageInfo_BookReference.Size(m)
}
func (m *BookReference) XXX_DiscardUnknown() {
	xxx_messageInfo_BookReference.DiscardUnknown(m)
}

var xxx_messageInfo_BookReference proto.InternalMessageInfo

func (m *BookReference) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func init() {
	proto.RegisterType((*Book)(nil), "Book")
	proto.RegisterType((*BookReference)(nil), "BookReference")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BookStoreClient is the client API for BookStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BookStoreClient interface {
	GetBook(ctx context.Context, in *BookReference, opts ...grpc.CallOption) (*Book, error)
}

type bookStoreClient struct {
	cc *grpc.ClientConn
}

func NewBookStoreClient(cc *grpc.ClientConn) BookStoreClient {
	return &bookStoreClient{cc}
}

func (c *bookStoreClient) GetBook(ctx context.Context, in *BookReference, opts ...grpc.CallOption) (*Book, error) {
	out := new(Book)
	err := c.cc.Invoke(ctx, "/BookStore/GetBook", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookStoreServer is the server API for BookStore service.
type BookStoreServer interface {
	GetBook(context.Context, *BookReference) (*Book, error)
}

func RegisterBookStoreServer(s *grpc.Server, srv BookStoreServer) {
	s.RegisterService(&_BookStore_serviceDesc, srv)
}

func _BookStore_GetBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookReference)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookStoreServer).GetBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/BookStore/GetBook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookStoreServer).GetBook(ctx, req.(*BookReference))
	}
	return interceptor(ctx, in, info, handler)
}

var _BookStore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "BookStore",
	HandlerType: (*BookStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBook",
			Handler:    _BookStore_GetBook_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "book_store.proto",
}

func init() { proto.RegisterFile("book_store.proto", fileDescriptor_book_store_1e7c586e95bda3a3) }

var fileDescriptor_book_store_1e7c586e95bda3a3 = []byte{
	// 151 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x48, 0xca, 0xcf, 0xcf,
	0x8e, 0x2f, 0x2e, 0xc9, 0x2f, 0x4a, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0xf2, 0xe0, 0x62,
	0x71, 0xca, 0xcf, 0xcf, 0x16, 0x12, 0xe1, 0x62, 0x2d, 0xc9, 0x2c, 0xc9, 0x49, 0x95, 0x60, 0x54,
	0x60, 0xd4, 0xe0, 0x0c, 0x82, 0x70, 0x84, 0xc4, 0xb8, 0xd8, 0x12, 0x4b, 0x4b, 0x32, 0xf2, 0x8b,
	0x24, 0x98, 0xc0, 0xc2, 0x50, 0x9e, 0x90, 0x10, 0x17, 0x4b, 0x65, 0x6a, 0x62, 0x91, 0x04, 0xb3,
	0x02, 0xa3, 0x06, 0x6b, 0x10, 0x98, 0xad, 0xa4, 0xca, 0xc5, 0x0b, 0x32, 0x29, 0x28, 0x35, 0x2d,
	0xb5, 0x28, 0x35, 0x2f, 0x39, 0x15, 0xbb, 0x91, 0x46, 0xfa, 0x5c, 0x9c, 0x20, 0x65, 0xc1, 0x20,
	0x37, 0x08, 0x29, 0x71, 0xb1, 0xbb, 0xa7, 0x96, 0x80, 0x1d, 0xc0, 0xa7, 0x87, 0xa2, 0x5b, 0x8a,
	0x15, 0xcc, 0x57, 0x62, 0x48, 0x62, 0x03, 0x3b, 0xd4, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x74,
	0x4d, 0x9c, 0xfb, 0xbc, 0x00, 0x00, 0x00,
}
