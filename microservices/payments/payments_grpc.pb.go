// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package paymentspb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PaymentsClient is the client API for Payments service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentsClient interface {
	CreateCharge(ctx context.Context, in *ChargeRequest, opts ...grpc.CallOption) (*ChargeResponse, error)
}

type paymentsClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentsClient(cc grpc.ClientConnInterface) PaymentsClient {
	return &paymentsClient{cc}
}

func (c *paymentsClient) CreateCharge(ctx context.Context, in *ChargeRequest, opts ...grpc.CallOption) (*ChargeResponse, error) {
	out := new(ChargeResponse)
	err := c.cc.Invoke(ctx, "/payments.Payments/CreateCharge", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentsServer is the server API for Payments service.
// All implementations must embed UnimplementedPaymentsServer
// for forward compatibility
type PaymentsServer interface {
	CreateCharge(context.Context, *ChargeRequest) (*ChargeResponse, error)
	mustEmbedUnimplementedPaymentsServer()
}

// UnimplementedPaymentsServer must be embedded to have forward compatible implementations.
type UnimplementedPaymentsServer struct {
}

func (UnimplementedPaymentsServer) CreateCharge(context.Context, *ChargeRequest) (*ChargeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCharge not implemented")
}
func (UnimplementedPaymentsServer) mustEmbedUnimplementedPaymentsServer() {}

// UnsafePaymentsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentsServer will
// result in compilation errors.
type UnsafePaymentsServer interface {
	mustEmbedUnimplementedPaymentsServer()
}

func RegisterPaymentsServer(s grpc.ServiceRegistrar, srv PaymentsServer) {
	s.RegisterService(&Payments_ServiceDesc, srv)
}

func _Payments_CreateCharge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChargeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServer).CreateCharge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payments.Payments/CreateCharge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServer).CreateCharge(ctx, req.(*ChargeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Payments_ServiceDesc is the grpc.ServiceDesc for Payments service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Payments_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "payments.Payments",
	HandlerType: (*PaymentsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCharge",
			Handler:    _Payments_CreateCharge_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "microservices/payments/payments.proto",
}
