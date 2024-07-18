// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.14.0
// source: bankaccount/protobuf/bankaccount.proto

package bankaccount

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PixService_Transfer_FullMethodName      = "/bankaccount.PixService/Transfer"
	PixService_DepositAmount_FullMethodName = "/bankaccount.PixService/DepositAmount"
	PixService_GetBalance_FullMethodName    = "/bankaccount.PixService/GetBalance"
)

// PixServiceClient is the client API for PixService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PixServiceClient interface {
	Transfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DepositAmount(ctx context.Context, in *DepositAmountRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error)
}

type pixServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPixServiceClient(cc grpc.ClientConnInterface) PixServiceClient {
	return &pixServiceClient{cc}
}

func (c *pixServiceClient) Transfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PixService_Transfer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pixServiceClient) DepositAmount(ctx context.Context, in *DepositAmountRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PixService_DepositAmount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pixServiceClient) GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error) {
	out := new(GetBalanceResponse)
	err := c.cc.Invoke(ctx, PixService_GetBalance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PixServiceServer is the server API for PixService service.
// All implementations must embed UnimplementedPixServiceServer
// for forward compatibility
type PixServiceServer interface {
	Transfer(context.Context, *TransferRequest) (*emptypb.Empty, error)
	DepositAmount(context.Context, *DepositAmountRequest) (*emptypb.Empty, error)
	GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
	mustEmbedUnimplementedPixServiceServer()
}

// UnimplementedPixServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPixServiceServer struct {
}

func (UnimplementedPixServiceServer) Transfer(context.Context, *TransferRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Transfer not implemented")
}
func (UnimplementedPixServiceServer) DepositAmount(context.Context, *DepositAmountRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DepositAmount not implemented")
}
func (UnimplementedPixServiceServer) GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalance not implemented")
}
func (UnimplementedPixServiceServer) mustEmbedUnimplementedPixServiceServer() {}

// UnsafePixServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PixServiceServer will
// result in compilation errors.
type UnsafePixServiceServer interface {
	mustEmbedUnimplementedPixServiceServer()
}

func RegisterPixServiceServer(s grpc.ServiceRegistrar, srv PixServiceServer) {
	s.RegisterService(&PixService_ServiceDesc, srv)
}

func _PixService_Transfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PixServiceServer).Transfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PixService_Transfer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PixServiceServer).Transfer(ctx, req.(*TransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PixService_DepositAmount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepositAmountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PixServiceServer).DepositAmount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PixService_DepositAmount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PixServiceServer).DepositAmount(ctx, req.(*DepositAmountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PixService_GetBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PixServiceServer).GetBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PixService_GetBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PixServiceServer).GetBalance(ctx, req.(*GetBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PixService_ServiceDesc is the grpc.ServiceDesc for PixService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PixService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bankaccount.PixService",
	HandlerType: (*PixServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Transfer",
			Handler:    _PixService_Transfer_Handler,
		},
		{
			MethodName: "DepositAmount",
			Handler:    _PixService_DepositAmount_Handler,
		},
		{
			MethodName: "GetBalance",
			Handler:    _PixService_GetBalance_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bankaccount/protobuf/bankaccount.proto",
}
