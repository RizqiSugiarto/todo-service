// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.21.12
// source: task/task_service.proto

package task

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
	TaskService_Create_FullMethodName         = "/proto.TaskService/Create"
	TaskService_Get_FullMethodName            = "/proto.TaskService/Get"
	TaskService_GetAllByUserID_FullMethodName = "/proto.TaskService/GetAllByUserID"
	TaskService_Update_FullMethodName         = "/proto.TaskService/Update"
	TaskService_Delete_FullMethodName         = "/proto.TaskService/Delete"
	TaskService_BatchUpdate_FullMethodName    = "/proto.TaskService/BatchUpdate"
)

// TaskServiceClient is the client API for TaskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TaskServiceClient interface {
	Create(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*TaskBaseResponse, error)
	Get(ctx context.Context, in *GetTaskByIDRequest, opts ...grpc.CallOption) (*GetTaskByIDResponse, error)
	GetAllByUserID(ctx context.Context, in *GetAllTaskByActivityIDRequest, opts ...grpc.CallOption) (*GetAllTaskByActivityIDResponse, error)
	Update(ctx context.Context, in *UpdateTaskByIDRequest, opts ...grpc.CallOption) (*TaskBaseResponse, error)
	Delete(ctx context.Context, in *DeleteTaskByIDRequest, opts ...grpc.CallOption) (*TaskBaseResponse, error)
	BatchUpdate(ctx context.Context, in *BatchUpdateTaskRequest, opts ...grpc.CallOption) (*TaskBaseResponse, error)
}

type taskServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskServiceClient(cc grpc.ClientConnInterface) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) Create(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*TaskBaseResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TaskBaseResponse)
	err := c.cc.Invoke(ctx, TaskService_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) Get(ctx context.Context, in *GetTaskByIDRequest, opts ...grpc.CallOption) (*GetTaskByIDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTaskByIDResponse)
	err := c.cc.Invoke(ctx, TaskService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) GetAllByUserID(ctx context.Context, in *GetAllTaskByActivityIDRequest, opts ...grpc.CallOption) (*GetAllTaskByActivityIDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllTaskByActivityIDResponse)
	err := c.cc.Invoke(ctx, TaskService_GetAllByUserID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) Update(ctx context.Context, in *UpdateTaskByIDRequest, opts ...grpc.CallOption) (*TaskBaseResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TaskBaseResponse)
	err := c.cc.Invoke(ctx, TaskService_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) Delete(ctx context.Context, in *DeleteTaskByIDRequest, opts ...grpc.CallOption) (*TaskBaseResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TaskBaseResponse)
	err := c.cc.Invoke(ctx, TaskService_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) BatchUpdate(ctx context.Context, in *BatchUpdateTaskRequest, opts ...grpc.CallOption) (*TaskBaseResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TaskBaseResponse)
	err := c.cc.Invoke(ctx, TaskService_BatchUpdate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskServiceServer is the server API for TaskService service.
// All implementations must embed UnimplementedTaskServiceServer
// for forward compatibility
type TaskServiceServer interface {
	Create(context.Context, *CreateTaskRequest) (*TaskBaseResponse, error)
	Get(context.Context, *GetTaskByIDRequest) (*GetTaskByIDResponse, error)
	GetAllByUserID(context.Context, *GetAllTaskByActivityIDRequest) (*GetAllTaskByActivityIDResponse, error)
	Update(context.Context, *UpdateTaskByIDRequest) (*TaskBaseResponse, error)
	Delete(context.Context, *DeleteTaskByIDRequest) (*TaskBaseResponse, error)
	BatchUpdate(context.Context, *BatchUpdateTaskRequest) (*TaskBaseResponse, error)
	mustEmbedUnimplementedTaskServiceServer()
}

// UnimplementedTaskServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTaskServiceServer struct {
}

func (UnimplementedTaskServiceServer) Create(context.Context, *CreateTaskRequest) (*TaskBaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedTaskServiceServer) Get(context.Context, *GetTaskByIDRequest) (*GetTaskByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedTaskServiceServer) GetAllByUserID(context.Context, *GetAllTaskByActivityIDRequest) (*GetAllTaskByActivityIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllByUserID not implemented")
}
func (UnimplementedTaskServiceServer) Update(context.Context, *UpdateTaskByIDRequest) (*TaskBaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedTaskServiceServer) Delete(context.Context, *DeleteTaskByIDRequest) (*TaskBaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedTaskServiceServer) BatchUpdate(context.Context, *BatchUpdateTaskRequest) (*TaskBaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchUpdate not implemented")
}
func (UnimplementedTaskServiceServer) mustEmbedUnimplementedTaskServiceServer() {}

// UnsafeTaskServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TaskServiceServer will
// result in compilation errors.
type UnsafeTaskServiceServer interface {
	mustEmbedUnimplementedTaskServiceServer()
}

func RegisterTaskServiceServer(s grpc.ServiceRegistrar, srv TaskServiceServer) {
	s.RegisterService(&TaskService_ServiceDesc, srv)
}

func _TaskService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).Create(ctx, req.(*CreateTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTaskByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).Get(ctx, req.(*GetTaskByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_GetAllByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllTaskByActivityIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).GetAllByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_GetAllByUserID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).GetAllByUserID(ctx, req.(*GetAllTaskByActivityIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTaskByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).Update(ctx, req.(*UpdateTaskByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTaskByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).Delete(ctx, req.(*DeleteTaskByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_BatchUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchUpdateTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).BatchUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_BatchUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).BatchUpdate(ctx, req.(*BatchUpdateTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TaskService_ServiceDesc is the grpc.ServiceDesc for TaskService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TaskService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _TaskService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _TaskService_Get_Handler,
		},
		{
			MethodName: "GetAllByUserID",
			Handler:    _TaskService_GetAllByUserID_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _TaskService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _TaskService_Delete_Handler,
		},
		{
			MethodName: "BatchUpdate",
			Handler:    _TaskService_BatchUpdate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "task/task_service.proto",
}
