// Copyright 2015 gRPC authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: protos/commandable.proto

package protos

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

const (
	Commandable_Invoke_FullMethodName = "/commandable.Commandable/invoke"
)

// CommandableClient is the client API for Commandable service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommandableClient interface {
	Invoke(ctx context.Context, in *InvokeRequest, opts ...grpc.CallOption) (*InvokeReply, error)
}

type commandableClient struct {
	cc grpc.ClientConnInterface
}

func NewCommandableClient(cc grpc.ClientConnInterface) CommandableClient {
	return &commandableClient{cc}
}

func (c *commandableClient) Invoke(ctx context.Context, in *InvokeRequest, opts ...grpc.CallOption) (*InvokeReply, error) {
	out := new(InvokeReply)
	err := c.cc.Invoke(ctx, Commandable_Invoke_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommandableServer is the server API for Commandable service.
// All implementations must embed UnimplementedCommandableServer
// for forward compatibility
type CommandableServer interface {
	Invoke(context.Context, *InvokeRequest) (*InvokeReply, error)
	mustEmbedUnimplementedCommandableServer()
}

// UnimplementedCommandableServer must be embedded to have forward compatible implementations.
type UnimplementedCommandableServer struct {
}

func (UnimplementedCommandableServer) Invoke(context.Context, *InvokeRequest) (*InvokeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Invoke not implemented")
}
func (UnimplementedCommandableServer) mustEmbedUnimplementedCommandableServer() {}

// UnsafeCommandableServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommandableServer will
// result in compilation errors.
type UnsafeCommandableServer interface {
	mustEmbedUnimplementedCommandableServer()
}

func RegisterCommandableServer(s grpc.ServiceRegistrar, srv CommandableServer) {
	s.RegisterService(&Commandable_ServiceDesc, srv)
}

func _Commandable_Invoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InvokeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandableServer).Invoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Commandable_Invoke_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandableServer).Invoke(ctx, req.(*InvokeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Commandable_ServiceDesc is the grpc.ServiceDesc for Commandable service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Commandable_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "commandable.Commandable",
	HandlerType: (*CommandableServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "invoke",
			Handler:    _Commandable_Invoke_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/commandable.proto",
}
