// Specifies the version of the protocol buffer syntax used

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: proto/protos/file.proto

// Specifies the package name for the generated code

package __

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
	FileService_AddFile_FullMethodName         = "/file.FileService/add_file"
	FileService_RemoveFile_FullMethodName      = "/file.FileService/remove_file"
	FileService_GetFileByFileid_FullMethodName = "/file.FileService/get_file_by_fileid"
	FileService_GetFileByUserid_FullMethodName = "/file.FileService/get_file_by_userid"
	FileService_GetFileByRoomid_FullMethodName = "/file.FileService/get_file_by_roomid"
	FileService_GetFiles_FullMethodName        = "/file.FileService/get_files"
)

// FileServiceClient is the client API for FileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileServiceClient interface {
	// Method to add a file
	AddFile(ctx context.Context, in *AddFileRequest, opts ...grpc.CallOption) (*AddFileResponse, error)
	// Method to remove a file
	RemoveFile(ctx context.Context, in *RemoveFileRequest, opts ...grpc.CallOption) (*Empty, error)
	// Method to get files by file ID
	GetFileByFileid(ctx context.Context, in *GetFileByFileIdRequest, opts ...grpc.CallOption) (*Files, error)
	// Method to get files by user ID
	GetFileByUserid(ctx context.Context, in *GetFileByUserIdRequest, opts ...grpc.CallOption) (*Files, error)
	// Method to get files by room ID
	GetFileByRoomid(ctx context.Context, in *GetFileByRoomIdRequest, opts ...grpc.CallOption) (*Files, error)
	GetFiles(ctx context.Context, in *Pagination, opts ...grpc.CallOption) (*Files, error)
}

type fileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileServiceClient(cc grpc.ClientConnInterface) FileServiceClient {
	return &fileServiceClient{cc}
}

func (c *fileServiceClient) AddFile(ctx context.Context, in *AddFileRequest, opts ...grpc.CallOption) (*AddFileResponse, error) {
	out := new(AddFileResponse)
	err := c.cc.Invoke(ctx, FileService_AddFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) RemoveFile(ctx context.Context, in *RemoveFileRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, FileService_RemoveFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) GetFileByFileid(ctx context.Context, in *GetFileByFileIdRequest, opts ...grpc.CallOption) (*Files, error) {
	out := new(Files)
	err := c.cc.Invoke(ctx, FileService_GetFileByFileid_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) GetFileByUserid(ctx context.Context, in *GetFileByUserIdRequest, opts ...grpc.CallOption) (*Files, error) {
	out := new(Files)
	err := c.cc.Invoke(ctx, FileService_GetFileByUserid_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) GetFileByRoomid(ctx context.Context, in *GetFileByRoomIdRequest, opts ...grpc.CallOption) (*Files, error) {
	out := new(Files)
	err := c.cc.Invoke(ctx, FileService_GetFileByRoomid_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) GetFiles(ctx context.Context, in *Pagination, opts ...grpc.CallOption) (*Files, error) {
	out := new(Files)
	err := c.cc.Invoke(ctx, FileService_GetFiles_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileServiceServer is the server API for FileService service.
// All implementations must embed UnimplementedFileServiceServer
// for forward compatibility
type FileServiceServer interface {
	// Method to add a file
	AddFile(context.Context, *AddFileRequest) (*AddFileResponse, error)
	// Method to remove a file
	RemoveFile(context.Context, *RemoveFileRequest) (*Empty, error)
	// Method to get files by file ID
	GetFileByFileid(context.Context, *GetFileByFileIdRequest) (*Files, error)
	// Method to get files by user ID
	GetFileByUserid(context.Context, *GetFileByUserIdRequest) (*Files, error)
	// Method to get files by room ID
	GetFileByRoomid(context.Context, *GetFileByRoomIdRequest) (*Files, error)
	GetFiles(context.Context, *Pagination) (*Files, error)
	mustEmbedUnimplementedFileServiceServer()
}

// UnimplementedFileServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFileServiceServer struct {
}

func (UnimplementedFileServiceServer) AddFile(context.Context, *AddFileRequest) (*AddFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddFile not implemented")
}
func (UnimplementedFileServiceServer) RemoveFile(context.Context, *RemoveFileRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveFile not implemented")
}
func (UnimplementedFileServiceServer) GetFileByFileid(context.Context, *GetFileByFileIdRequest) (*Files, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileByFileid not implemented")
}
func (UnimplementedFileServiceServer) GetFileByUserid(context.Context, *GetFileByUserIdRequest) (*Files, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileByUserid not implemented")
}
func (UnimplementedFileServiceServer) GetFileByRoomid(context.Context, *GetFileByRoomIdRequest) (*Files, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileByRoomid not implemented")
}
func (UnimplementedFileServiceServer) GetFiles(context.Context, *Pagination) (*Files, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFiles not implemented")
}
func (UnimplementedFileServiceServer) mustEmbedUnimplementedFileServiceServer() {}

// UnsafeFileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileServiceServer will
// result in compilation errors.
type UnsafeFileServiceServer interface {
	mustEmbedUnimplementedFileServiceServer()
}

func RegisterFileServiceServer(s grpc.ServiceRegistrar, srv FileServiceServer) {
	s.RegisterService(&FileService_ServiceDesc, srv)
}

func _FileService_AddFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).AddFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_AddFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).AddFile(ctx, req.(*AddFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_RemoveFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).RemoveFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_RemoveFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).RemoveFile(ctx, req.(*RemoveFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_GetFileByFileid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileByFileIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).GetFileByFileid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_GetFileByFileid_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).GetFileByFileid(ctx, req.(*GetFileByFileIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_GetFileByUserid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileByUserIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).GetFileByUserid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_GetFileByUserid_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).GetFileByUserid(ctx, req.(*GetFileByUserIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_GetFileByRoomid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileByRoomIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).GetFileByRoomid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_GetFileByRoomid_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).GetFileByRoomid(ctx, req.(*GetFileByRoomIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_GetFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Pagination)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).GetFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_GetFiles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).GetFiles(ctx, req.(*Pagination))
	}
	return interceptor(ctx, in, info, handler)
}

// FileService_ServiceDesc is the grpc.ServiceDesc for FileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "file.FileService",
	HandlerType: (*FileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "add_file",
			Handler:    _FileService_AddFile_Handler,
		},
		{
			MethodName: "remove_file",
			Handler:    _FileService_RemoveFile_Handler,
		},
		{
			MethodName: "get_file_by_fileid",
			Handler:    _FileService_GetFileByFileid_Handler,
		},
		{
			MethodName: "get_file_by_userid",
			Handler:    _FileService_GetFileByUserid_Handler,
		},
		{
			MethodName: "get_file_by_roomid",
			Handler:    _FileService_GetFileByRoomid_Handler,
		},
		{
			MethodName: "get_files",
			Handler:    _FileService_GetFiles_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/protos/file.proto",
}
