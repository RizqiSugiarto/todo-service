package handler

import (
	"context"

	"github.com/digisata/invitation-service/internal/entity"
)

type (
	InvitationUseCase interface {
		CreateInvitation(ctx context.Context, req entity.CreateInvitationRequest) error
		UpdateInvitation(ctx context.Context, req entity.UpdateInvitationRequest) error
		GetInvitation(ctx context.Context, id string) (entity.Invitation, error)
		GetAllInvitationByUserID(ctx context.Context, req entity.GetAllInvitationRequest) ([]entity.Invitation, entity.Paging, error)
		DeleteInvitation(ctx context.Context, id string) error
		GetInvitationLink(ctx context.Context, id string) (string, error)
		VerifyInvitation(ctx context.Context, id string) (entity.VerifyInvitationResponse, error)
	}

	InvitationLabelUseCase interface {
		CreateInvitationLabel(ctx context.Context, req entity.CreateInvitationLabelRequest) error
		UpdateInvitationLabel(ctx context.Context, req entity.UpdateInvitationLabelRequest) error
		GetInvitationLabel(ctx context.Context, id string) (entity.InvitationLabel, error)
		GetAllInvitationLabelByUserID(ctx context.Context, req entity.GetAllInvitationLabelRequest) ([]entity.InvitationLabel, entity.Paging, error)
		DeleteInvitationLabel(ctx context.Context, id string) error
	}

	InvitationCategoryUseCase interface {
		CreateInvitationCategory(ctx context.Context, req entity.CreateInvitationCategoryRequest) error
		UpdateInvitationCategory(ctx context.Context, req entity.UpdateInvitationCategoryRequest) error
		GetInvitationCategory(ctx context.Context, id string) (entity.InvitationCategory, error)
		GetAllInvitationCategory(ctx context.Context, req entity.GetAllInvitationCategoryRequest) ([]entity.InvitationCategory, entity.Paging, error)
		DeleteInvitationCategory(ctx context.Context, id string) error
	}
)
