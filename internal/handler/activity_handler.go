package handler

import (
	"context"

	"github.com/digisata/todo-service/internal/entity"
	activityPB "github.com/digisata/todo-service/stubs/activity"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ActivityHandler struct {
	activityPB.UnimplementedActivityServiceServer
	activityUseCase ActivityUseCase
}

func NewActivity(activityUseCase ActivityUseCase) *ActivityHandler {
	return &ActivityHandler{
		activityUseCase: activityUseCase,
	}
}

func (h *ActivityHandler) Create(ctx context.Context, req *activityPB.CreateActivityRequest) (*activityPB.ActivityBaseResponse, error) {
	payload := entity.CreateActivityRequest{
		Title: req.GetTitle(),
		Type:  req.GetType(),
	}

	// if payload.Type == "text" {
	// 	h.acitivityTextUseCase.CreateText(ctx, entity.CreateTextRequest{
	// 		ActivityID: pa,
	// 	})
	// }

	err := h.activityUseCase.CreateActivity(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &activityPB.ActivityBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (g *ActivityHandler) Update(ctx context.Context, req *activityPB.UpdateActivityByIDRequest) (*activityPB.ActivityBaseResponse, error) {
	payload := entity.UpdateActivityRequest{
		ID:    req.GetId(),
		Title: req.GetTitle(),
		Type:  req.GetType(),
	}

	err := g.activityUseCase.UpdateActivity(ctx, payload)
	if err != nil && err.Error() == "data not found" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &activityPB.ActivityBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (h *ActivityHandler) Get(ctx context.Context, req *activityPB.GetActivityByIDRequest) (*activityPB.ActivityBaseResponse, error) {
	data, err := h.activityUseCase.GetActivity(ctx, req.GetId())
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	getActivityByIDResponse := &activityPB.GetActivityByIDResponse{
		Id:        data.ID,
		Title:     data.Title,
		Type:      data.Type,
		CreatedAt: timestamppb.New(data.CreatedAt),
		UpdatedAt: timestamppb.New(data.UpdatedAt),
	}

	anyData, err := anypb.New(getActivityByIDResponse)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &activityPB.ActivityBaseResponse{
		Message: "Success",
		Data:    anyData,
	}

	return res, nil
}

func (g *ActivityHandler) GetAll(ctx context.Context, req *activityPB.GetAllActivityRequest) (*activityPB.GetAllActivityResponse, error) {
	payload := entity.GetAllActivityRequest{
		Search: req.Search,
		Page:   req.Page,
		Limit:  req.Limit,
	}

	data, paging, err := g.activityUseCase.GetAllActivity(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &activityPB.GetAllActivityResponse{
		Message: "Success",
		Data:    []*activityPB.GetActivityByIDResponse{},
		Paging: &activityPB.ActivityPaging{
			CurrentPage: paging.CurrentPage,
			TotalPage:   paging.TotalPage,
			Count:       paging.Count,
		},
	}

	for _, activity := range data {
		data := &activityPB.GetActivityByIDResponse{
			Id:        activity.ID,
			Title:     activity.Title,
			Type:      activity.Type,
			CreatedAt: timestamppb.New(activity.CreatedAt),
			UpdatedAt: timestamppb.New(activity.UpdatedAt),
		}

		res.Data = append(res.Data, data)
	}

	return res, nil
}

func (g *ActivityHandler) Delete(ctx context.Context, req *activityPB.DeleteActivityByIDRequest) (*activityPB.ActivityBaseResponse, error) {
	err := g.activityUseCase.DeleteActivity(ctx, req.GetId())
	if err != nil && err.Error() == "data not found" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &activityPB.ActivityBaseResponse{
		Message: "Success",
	}

	return res, nil
}
