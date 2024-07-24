package handler

import (
	"context"

	"github.com/digisata/invitation-service/internal/entity"
	invitationLabelPB "github.com/digisata/invitation-service/stubs/invitation-label"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InvitationLabelHandler struct {
	invitationLabelPB.UnimplementedInvitationLabelServiceServer
	invitationLabelUseCase InvitationLabelUseCase
}

func NewInvitationLabel(invitationLabelUseCase InvitationLabelUseCase) *InvitationLabelHandler {
	return &InvitationLabelHandler{
		invitationLabelUseCase: invitationLabelUseCase,
	}
}

func (h *InvitationLabelHandler) CreateInvitationLabel(ctx context.Context, req *invitationLabelPB.CreateInvitationLabelRequest) (*invitationLabelPB.InvitationLabelBaseResponse, error) {
	payload := entity.CreateInvitationLabelRequest{
		UserID: req.GetUserId(),
		Name:   req.GetName(),
	}

	err := h.invitationLabelUseCase.CreateInvitationLabel(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationLabelPB.InvitationLabelBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (g *InvitationLabelHandler) UpdateInvitationLabel(ctx context.Context, req *invitationLabelPB.UpdateInvitationLabelByIDRequest) (*invitationLabelPB.InvitationLabelBaseResponse, error) {
	payload := entity.UpdateInvitationLabelRequest{
		ID:   req.GetId(),
		Name: req.GetName(),
	}

	err := g.invitationLabelUseCase.UpdateInvitationLabel(ctx, payload)
	if err != nil && err.Error() == "data not found" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationLabelPB.InvitationLabelBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (h *InvitationLabelHandler) GetInvitationLabel(ctx context.Context, req *invitationLabelPB.GetInvitationLabelByIDRequest) (*invitationLabelPB.GetInvitationLabelByIDResponse, error) {
	data, err := h.invitationLabelUseCase.GetInvitationLabel(ctx, req.GetId())
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationLabelPB.GetInvitationLabelByIDResponse{
		Id:        data.ID,
		UserId:    data.UserID,
		Name:      data.Name,
		CreatedAt: timestamppb.New(data.CreatedAt),
		UpdatedAt: timestamppb.New(data.UpdatedAt),
	}

	return res, nil
}

func (g *InvitationLabelHandler) GetAllInvitationLabelByUserID(ctx context.Context, req *invitationLabelPB.GetAllInvitationLabelByUserIDRequest) (*invitationLabelPB.GetAllInvitationLabelByUserIDResponse, error) {
	payload := entity.GetAllInvitationLabelRequest{
		UserID: req.GetUserId(),
		Search: req.Search,
		Page:   req.Page,
		Limit:  req.Limit,
	}

	data, paging, err := g.invitationLabelUseCase.GetAllInvitationLabelByUserID(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationLabelPB.GetAllInvitationLabelByUserIDResponse{
		Message:          "Success",
		InvitationLabels: []*invitationLabelPB.GetInvitationLabelByIDResponse{},
		Paging: &invitationLabelPB.InvitationLabelPaging{
			CurrentPage: paging.CurrentPage,
			TotalPage:   paging.TotalPage,
			Count:       paging.Count,
		},
	}
	for _, invitationLabel := range data {
		data := &invitationLabelPB.GetInvitationLabelByIDResponse{
			Id:        invitationLabel.ID,
			UserId:    invitationLabel.UserID,
			Name:      invitationLabel.Name,
			CreatedAt: timestamppb.New(invitationLabel.CreatedAt),
			UpdatedAt: timestamppb.New(invitationLabel.UpdatedAt),
		}

		res.InvitationLabels = append(res.InvitationLabels, data)
	}

	return res, nil
}

func (g *InvitationLabelHandler) DeleteInvitationLabel(ctx context.Context, req *invitationLabelPB.DeleteInvitationLabelByIDRequest) (*invitationLabelPB.InvitationLabelBaseResponse, error) {
	err := g.invitationLabelUseCase.DeleteInvitationLabel(ctx, req.GetId())
	if err != nil && err.Error() == "data not found" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationLabelPB.InvitationLabelBaseResponse{
		Message: "Success",
	}

	return res, nil
}
