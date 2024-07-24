package handler

import (
	"context"

	"github.com/digisata/invitation-service/internal/entity"
	invitationPB "github.com/digisata/invitation-service/stubs/invitation"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InvitationHandler struct {
	invitationPB.UnimplementedInvitationServiceServer
	invitationUseCase InvitationUseCase
}

func NewInvitation(invitationUseCase InvitationUseCase) *InvitationHandler {
	return &InvitationHandler{
		invitationUseCase: invitationUseCase,
	}
}

func (h *InvitationHandler) Create(ctx context.Context, req *invitationPB.CreateInvitationRequest) (*invitationPB.InvitationBaseResponse, error) {
	payload := entity.CreateInvitationRequest{
		UserID:               req.GetUserId(),
		Name:                 req.GetName(),
		InvitationLabelID:    req.InvitationLabelId,
		InvitationCategoryID: req.InvitationCategoryId,
	}

	err := h.invitationUseCase.CreateInvitation(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationPB.InvitationBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (g *InvitationHandler) Update(ctx context.Context, req *invitationPB.UpdateInvitationByIDRequest) (*invitationPB.InvitationBaseResponse, error) {
	payload := entity.UpdateInvitationRequest{
		ID:                   req.GetId(),
		Name:                 req.Name,
		InvitationLabelID:    req.InvitationLabelId,
		InvitationCategoryID: req.InvitationCategoryId,
	}

	err := g.invitationUseCase.UpdateInvitation(ctx, payload)
	if err != nil && err.Error() == "data not found" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationPB.InvitationBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (h *InvitationHandler) Get(ctx context.Context, req *invitationPB.GetInvitationByIDRequest) (*invitationPB.GetInvitationByIDResponse, error) {
	data, err := h.invitationUseCase.GetInvitation(ctx, req.GetId())
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationPB.GetInvitationByIDResponse{
		Id:          data.ID,
		UserId:      data.UserID,
		Name:        data.Name,
		IsOpen:      data.IsOpen,
		IsComing:    data.IsComing,
		IsSendMoney: data.IsSendMoney,
		IsSendGift:  data.IsSendGift,
		IsCheckIn:   data.IsCheckIn,
		CreatedAt:   timestamppb.New(data.CreatedAt),
		UpdatedAt:   timestamppb.New(data.UpdatedAt),
	}

	if data.InvitationLabelID != nil {
		res.InvitationLabel = &invitationPB.GetInvitationLabelResponse{
			Id:        *data.InvitationLabelID,
			UserId:    *data.InvitationLabelUserID,
			Name:      *data.InvitationLabelName,
			CreatedAt: timestamppb.New(*data.InvitationLabelCreatedAt),
			UpdatedAt: timestamppb.New(*data.InvitationLabelUpdatedAt),
		}
	}

	if data.InvitationCategoryID != nil {
		res.InvitationCategory = &invitationPB.GetInvitationCategoryResponse{
			Id:        *data.InvitationCategoryID,
			Name:      *data.InvitationCategoryName,
			CreatedAt: timestamppb.New(*data.InvitationCategoryCreatedAt),
			UpdatedAt: timestamppb.New(*data.InvitationCategoryUpdatedAt),
		}
	}

	if data.OpenAt != nil {
		res.OpenAt = timestamppb.New(*data.OpenAt)
	}

	if data.CheckInAt != nil {
		res.CheckInAt = timestamppb.New(*data.CheckInAt)
	}

	return res, nil
}

func (g *InvitationHandler) GetAllByUserID(ctx context.Context, req *invitationPB.GetAllInvitationByUserIDRequest) (*invitationPB.GetAllInvitationByUserIDResponse, error) {
	payload := entity.GetAllInvitationRequest{
		UserID:               req.GetUserId(),
		Search:               req.Search,
		Page:                 req.Page,
		Limit:                req.Limit,
		IsOpen:               req.IsOpen,
		IsComing:             req.IsComing,
		IsSendMoney:          req.IsSendMoney,
		IsSendGift:           req.IsSendGift,
		IsCheckIn:            req.IsCheckIn,
		InvitationLabels:     req.InvitationLabels,
		InvitationCategories: req.InvitationCategories,
	}

	data, paging, err := g.invitationUseCase.GetAllInvitationByUserID(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationPB.GetAllInvitationByUserIDResponse{
		Message:     "Success",
		Invitations: []*invitationPB.GetInvitationByIDResponse{},
		Paging: &invitationPB.InvitationPaging{
			CurrentPage: paging.CurrentPage,
			TotalPage:   paging.TotalPage,
			Count:       paging.Count,
		},
	}
	for _, invitation := range data {
		data := &invitationPB.GetInvitationByIDResponse{
			Id:          invitation.ID,
			UserId:      invitation.UserID,
			Name:        invitation.Name,
			IsOpen:      invitation.IsOpen,
			IsComing:    invitation.IsComing,
			IsSendMoney: invitation.IsSendMoney,
			IsSendGift:  invitation.IsSendGift,
			IsCheckIn:   invitation.IsCheckIn,
			CreatedAt:   timestamppb.New(invitation.CreatedAt),
			UpdatedAt:   timestamppb.New(invitation.UpdatedAt),
		}

		if invitation.InvitationLabelID != nil {
			data.InvitationLabel = &invitationPB.GetInvitationLabelResponse{
				Id:        *invitation.InvitationLabelID,
				UserId:    *invitation.InvitationLabelUserID,
				Name:      *invitation.InvitationLabelName,
				CreatedAt: timestamppb.New(*invitation.InvitationLabelCreatedAt),
				UpdatedAt: timestamppb.New(*invitation.InvitationLabelUpdatedAt),
			}
		}

		if invitation.InvitationCategoryID != nil {
			data.InvitationCategory = &invitationPB.GetInvitationCategoryResponse{
				Id:        *invitation.InvitationCategoryID,
				Name:      *invitation.InvitationCategoryName,
				CreatedAt: timestamppb.New(*invitation.InvitationCategoryCreatedAt),
				UpdatedAt: timestamppb.New(*invitation.InvitationCategoryUpdatedAt),
			}
		}

		if invitation.OpenAt != nil {
			data.OpenAt = timestamppb.New(*invitation.OpenAt)
		}

		if invitation.CheckInAt != nil {
			data.CheckInAt = timestamppb.New(*invitation.CheckInAt)
		}

		res.Invitations = append(res.Invitations, data)
	}

	return res, nil
}

func (g *InvitationHandler) Delete(ctx context.Context, req *invitationPB.DeleteInvitationByIDRequest) (*invitationPB.InvitationBaseResponse, error) {
	err := g.invitationUseCase.DeleteInvitation(ctx, req.GetId())
	if err != nil && err.Error() == "data not found" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationPB.InvitationBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (g *InvitationHandler) GetLink(ctx context.Context, req *invitationPB.GetInvitationLinkByIDRequest) (*invitationPB.GetInvitationLinkByIDResponse, error) {
	data, err := g.invitationUseCase.GetInvitationLink(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationPB.GetInvitationLinkByIDResponse{
		Link: data,
	}

	return res, nil
}

func (g *InvitationHandler) Verify(ctx context.Context, req *invitationPB.VerifyInvitationRequest) (*invitationPB.VerifyInvitationResponse, error) {
	data, err := g.invitationUseCase.VerifyInvitation(ctx, req.GetId())
	if err != nil && err.Error() == "the invitation is not in the list" {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &invitationPB.VerifyInvitationResponse{
		UserId:       data.UserID,
		InvitationId: data.InvitationID,
	}

	return res, nil
}
