// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.12.4
// source: cache.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	CacheService_SetData_FullMethodName  = "/proto.CacheService/SetData"
	CacheService_GetData_FullMethodName  = "/proto.CacheService/GetData"
	CacheService_SyncData_FullMethodName = "/proto.CacheService/SyncData"
)

// CacheServiceClient is the client API for CacheService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CacheServiceClient interface {
	SetData(ctx context.Context, in *DataRequest, opts ...grpc.CallOption) (*DataResponse, error)
	GetData(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*DataResponse, error)
	SyncData(ctx context.Context, in *DataRequest, opts ...grpc.CallOption) (*DataResponse, error)
}

type cacheServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCacheServiceClient(cc grpc.ClientConnInterface) CacheServiceClient {
	return &cacheServiceClient{cc}
}

func (c *cacheServiceClient) SetData(ctx context.Context, in *DataRequest, opts ...grpc.CallOption) (*DataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DataResponse)
	err := c.cc.Invoke(ctx, CacheService_SetData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cacheServiceClient) GetData(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*DataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DataResponse)
	err := c.cc.Invoke(ctx, CacheService_GetData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cacheServiceClient) SyncData(ctx context.Context, in *DataRequest, opts ...grpc.CallOption) (*DataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DataResponse)
	err := c.cc.Invoke(ctx, CacheService_SyncData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CacheServiceServer is the server API for CacheService service.
// All implementations must embed UnimplementedCacheServiceServer
// for forward compatibility
type CacheServiceServer interface {
	SetData(context.Context, *DataRequest) (*DataResponse, error)
	GetData(context.Context, *KeyRequest) (*DataResponse, error)
	SyncData(context.Context, *DataRequest) (*DataResponse, error)
	mustEmbedUnimplementedCacheServiceServer()
}

// UnimplementedCacheServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCacheServiceServer struct {
}

func (UnimplementedCacheServiceServer) SetData(context.Context, *DataRequest) (*DataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetData not implemented")
}
func (UnimplementedCacheServiceServer) GetData(context.Context, *KeyRequest) (*DataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetData not implemented")
}
func (UnimplementedCacheServiceServer) SyncData(context.Context, *DataRequest) (*DataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncData not implemented")
}
func (UnimplementedCacheServiceServer) mustEmbedUnimplementedCacheServiceServer() {}

// UnsafeCacheServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CacheServiceServer will
// result in compilation errors.
type UnsafeCacheServiceServer interface {
	mustEmbedUnimplementedCacheServiceServer()
}

func RegisterCacheServiceServer(s grpc.ServiceRegistrar, srv CacheServiceServer) {
	s.RegisterService(&CacheService_ServiceDesc, srv)
}

func _CacheService_SetData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacheServiceServer).SetData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CacheService_SetData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacheServiceServer).SetData(ctx, req.(*DataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CacheService_GetData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacheServiceServer).GetData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CacheService_GetData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacheServiceServer).GetData(ctx, req.(*KeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CacheService_SyncData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacheServiceServer).SyncData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CacheService_SyncData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacheServiceServer).SyncData(ctx, req.(*DataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CacheService_ServiceDesc is the grpc.ServiceDesc for CacheService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CacheService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.CacheService",
	HandlerType: (*CacheServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetData",
			Handler:    _CacheService_SetData_Handler,
		},
		{
			MethodName: "GetData",
			Handler:    _CacheService_GetData_Handler,
		},
		{
			MethodName: "SyncData",
			Handler:    _CacheService_SyncData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cache.proto",
}
