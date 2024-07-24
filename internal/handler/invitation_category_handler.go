package handler

import (
	"context"

	"github.com/digisata/invitation-service/internal/entity"
	invitationCategoryPB "github.com/digisata/invitation-service/stubs/invitation-category"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InvitationCategoryHandler struct {
	invitationCategoryPB.UnimplementedInvitationCategoryServiceServer
	invitationCategoryUseCase InvitationCategoryUseCase
}

func NewInvitationCategory(invitationCategoryUseCase InvitationCategoryUseCase) *InvitationCategoryHandler {
	return &InvitationCategoryHandler{
		invitationCategoryUseCase: invitationCategoryUseCase,
	}
}

func (h *InvitationCategoryHandler) CreateInvitationCategory(ctx context.Context, req *invitationCategoryPB.CreateInvitationCategoryRequest) (*invitationCategoryPB.InvitationCategoryBaseResponse, error) {
	payload := entity.CreateInvitationCategoryRequest{
		Name: req.GetName(),
	}

	err := h.invitationCategoryUseCase.CreateInvitationCategory(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationCategoryPB.InvitationCategoryBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (g *InvitationCategoryHandler) UpdateInvitationCategory(ctx context.Context, req *invitationCategoryPB.UpdateInvitationCategoryByIDRequest) (*invitationCategoryPB.InvitationCategoryBaseResponse, error) {
	payload := entity.UpdateInvitationCategoryRequest{
		ID:   req.GetId(),
		Name: req.GetName(),
	}

	err := g.invitationCategoryUseCase.UpdateInvitationCategory(ctx, payload)
	if err != nil && err.Error() == "data not found" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationCategoryPB.InvitationCategoryBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (h *InvitationCategoryHandler) GetInvitationCategory(ctx context.Context, req *invitationCategoryPB.GetInvitationCategoryByIDRequest) (*invitationCategoryPB.InvitationCategoryBaseResponse, error) {
	data, err := h.invitationCategoryUseCase.GetInvitationCategory(ctx, req.GetId())
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	getInvitationCategoryByIDResponse := &invitationCategoryPB.GetInvitationCategoryByIDResponse{
		Id:        data.ID,
		Name:      data.Name,
		CreatedAt: timestamppb.New(data.CreatedAt),
		UpdatedAt: timestamppb.New(data.UpdatedAt),
	}

	anyData, err := anypb.New(getInvitationCategoryByIDResponse)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &invitationCategoryPB.InvitationCategoryBaseResponse{
		Message: "Success",
		Data:    anyData,
	}

	return res, nil
}

func (g *InvitationCategoryHandler) GetAllInvitationCategory(ctx context.Context, req *invitationCategoryPB.GetAllInvitationCategoryRequest) (*invitationCategoryPB.GetAllInvitationCategoryResponse, error) {
	payload := entity.GetAllInvitationCategoryRequest{
		Search: req.Search,
		Page:   req.Page,
		Limit:  req.Limit,
	}

	data, paging, err := g.invitationCategoryUseCase.GetAllInvitationCategory(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationCategoryPB.GetAllInvitationCategoryResponse{
		Message: "Success",
		Data:    []*invitationCategoryPB.GetInvitationCategoryByIDResponse{},
		Paging: &invitationCategoryPB.InvitationCategoryPaging{
			CurrentPage: paging.CurrentPage,
			TotalPage:   paging.TotalPage,
			Count:       paging.Count,
		},
	}

	for _, invitationCategory := range data {
		data := &invitationCategoryPB.GetInvitationCategoryByIDResponse{
			Id:        invitationCategory.ID,
			Name:      invitationCategory.Name,
			CreatedAt: timestamppb.New(invitationCategory.CreatedAt),
			UpdatedAt: timestamppb.New(invitationCategory.UpdatedAt),
		}

		res.Data = append(res.Data, data)
	}

	return res, nil
}

func (g *InvitationCategoryHandler) DeleteInvitationCategory(ctx context.Context, req *invitationCategoryPB.DeleteInvitationCategoryByIDRequest) (*invitationCategoryPB.InvitationCategoryBaseResponse, error) {
	err := g.invitationCategoryUseCase.DeleteInvitationCategory(ctx, req.GetId())
	if err != nil && err.Error() == "data not found" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationCategoryPB.InvitationCategoryBaseResponse{
		Message: "Success",
	}

	return res, nil
}
