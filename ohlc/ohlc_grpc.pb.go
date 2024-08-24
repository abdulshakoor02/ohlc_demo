// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: ohlc.proto

package ohlc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	OHLCService_StreamOHLCData_FullMethodName = "/ohlc.OHLCService/StreamOHLCData"
)

// OHLCServiceClient is the client API for OHLCService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OHLCServiceClient interface {
	StreamOHLCData(ctx context.Context, in *OHLCrequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[OHLC], error)
}

type oHLCServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOHLCServiceClient(cc grpc.ClientConnInterface) OHLCServiceClient {
	return &oHLCServiceClient{cc}
}

func (c *oHLCServiceClient) StreamOHLCData(ctx context.Context, in *OHLCrequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[OHLC], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &OHLCService_ServiceDesc.Streams[0], OHLCService_StreamOHLCData_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[OHLCrequest, OHLC]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type OHLCService_StreamOHLCDataClient = grpc.ServerStreamingClient[OHLC]

// OHLCServiceServer is the server API for OHLCService service.
// All implementations must embed UnimplementedOHLCServiceServer
// for forward compatibility.
type OHLCServiceServer interface {
	StreamOHLCData(*OHLCrequest, grpc.ServerStreamingServer[OHLC]) error
	mustEmbedUnimplementedOHLCServiceServer()
}

// UnimplementedOHLCServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedOHLCServiceServer struct{}

func (UnimplementedOHLCServiceServer) StreamOHLCData(*OHLCrequest, grpc.ServerStreamingServer[OHLC]) error {
	return status.Errorf(codes.Unimplemented, "method StreamOHLCData not implemented")
}
func (UnimplementedOHLCServiceServer) mustEmbedUnimplementedOHLCServiceServer() {}
func (UnimplementedOHLCServiceServer) testEmbeddedByValue()                     {}

// UnsafeOHLCServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OHLCServiceServer will
// result in compilation errors.
type UnsafeOHLCServiceServer interface {
	mustEmbedUnimplementedOHLCServiceServer()
}

func RegisterOHLCServiceServer(s grpc.ServiceRegistrar, srv OHLCServiceServer) {
	// If the following call pancis, it indicates UnimplementedOHLCServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&OHLCService_ServiceDesc, srv)
}

func _OHLCService_StreamOHLCData_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(OHLCrequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OHLCServiceServer).StreamOHLCData(m, &grpc.GenericServerStream[OHLCrequest, OHLC]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type OHLCService_StreamOHLCDataServer = grpc.ServerStreamingServer[OHLC]

// OHLCService_ServiceDesc is the grpc.ServiceDesc for OHLCService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OHLCService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ohlc.OHLCService",
	HandlerType: (*OHLCServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamOHLCData",
			Handler:       _OHLCService_StreamOHLCData_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "ohlc.proto",
}
