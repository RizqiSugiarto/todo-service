package handler

import (
	"context"

	"github.com/digisata/todo-service/internal/entity"
	taskPB "github.com/digisata/todo-service/stubs/task"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskHandler struct {
	taskPB.UnimplementedTaskServiceServer
	taskUseCase TaskUseCase
}

func NewTask(taskUseCase TaskUseCase) *TaskHandler {
	return &TaskHandler{
		taskUseCase: taskUseCase,
	}
}

func (h *TaskHandler) Create(ctx context.Context, req *taskPB.CreateTaskRequest) (*taskPB.TaskBaseResponse, error) {
	payload := entity.CreateTaskRequest{
		ActivityID: req.GetActivityId(),
		Title:      req.GetTitle(),
		IsActive:   req.IsActive,
		Priority:   int(req.GetPriority()),
	}

	err := h.taskUseCase.CreateTask(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &taskPB.TaskBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (g *TaskHandler) Update(ctx context.Context, req *taskPB.UpdateTaskByIDRequest) (*taskPB.TaskBaseResponse, error) {
	payload := entity.UpdateTaskRequest{
		ID:       req.GetId(),
		Title:    req.Title,
		IsActive: req.IsActive,
	}

	if req.Priority != nil {
		priority := int(*req.Priority)
		payload.Priority = &priority
	}

	err := g.taskUseCase.UpdateTask(ctx, payload)
	if err != nil && err.Error() == "data not found" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &taskPB.TaskBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (g *TaskHandler) BatchUpdate(ctx context.Context, req *taskPB.BatchUpdateTaskRequest) (*taskPB.TaskBaseResponse, error) {
	var payload []entity.UpdateTaskRequest

	for _, task := range req.Tasks {
		taskPayload := entity.UpdateTaskRequest{
			ID:       task.GetId(),
			Title:    task.Title,
			IsActive: task.IsActive,
		}

		if task.Priority != nil {
			priority := int(*task.Priority)
			taskPayload.Priority = &priority
		}

		if task.Order != nil {
			order := int(*task.Order)
			taskPayload.Order = &order
		}

		payload = append(payload, taskPayload)
	}

	err := g.taskUseCase.BatchUpdateTask(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &taskPB.TaskBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (h *TaskHandler) Get(ctx context.Context, req *taskPB.GetTaskByIDRequest) (*taskPB.GetTaskByIDResponse, error) {
	data, err := h.taskUseCase.GetTask(ctx, req.GetId())
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &taskPB.GetTaskByIDResponse{
		Id:         data.ID,
		ActivityId: data.ActivityID,
		Title:      data.Title,
		IsActive:   data.IsActive,
		Priority:   int32(data.Priority),
		Order:      int32(data.Order),
		CreatedAt:  timestamppb.New(data.CreatedAt),
		UpdatedAt:  timestamppb.New(data.UpdatedAt),
	}

	return res, nil
}

func (g *TaskHandler) GetAllByUserID(ctx context.Context, req *taskPB.GetAllTaskByActivityIDRequest) (*taskPB.GetAllTaskByActivityIDResponse, error) {
	payload := entity.GetAllTaskRequest{
		ActivityID:   req.GetActivityId(),
		Search:       req.Search,
		Page:         req.Page,
		Limit:        req.Limit,
		IsActive:     req.IsActive,
		IsNewest:     req.IsNewest,
		IsOldest:     req.IsOldest,
		IsAscending:  req.IsAscending,
		IsDescending: req.IsDescending,
	}

	if req.Priority != nil {
		priority := int(*req.Priority)
		payload.Priority = &priority
	}

	data, paging, err := g.taskUseCase.GetAllTaskByActivityID(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &taskPB.GetAllTaskByActivityIDResponse{
		Message: "Success",
		Tasks:   []*taskPB.GetTaskByIDResponse{},
		Paging: &taskPB.TaskPaging{
			CurrentPage: paging.CurrentPage,
			TotalPage:   paging.TotalPage,
			Count:       paging.Count,
		},
	}
	for _, task := range data {
		data := &taskPB.GetTaskByIDResponse{
			Id:         task.ID,
			ActivityId: task.ActivityID,
			Title:      task.Title,
			IsActive:   task.IsActive,
			Priority:   int32(task.Priority),
			Order:      int32(task.Order),
			CreatedAt:  timestamppb.New(task.CreatedAt),
			UpdatedAt:  timestamppb.New(task.UpdatedAt),
		}

		res.Tasks = append(res.Tasks, data)
	}

	return res, nil
}

func (g *TaskHandler) Delete(ctx context.Context, req *taskPB.DeleteTaskByIDRequest) (*taskPB.TaskBaseResponse, error) {
	err := g.taskUseCase.DeleteTask(ctx, req.GetId())
	if err != nil && err.Error() == "data not found" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &taskPB.TaskBaseResponse{
		Message: "Success",
	}

	return res, nil
}
