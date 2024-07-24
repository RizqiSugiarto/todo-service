package usecase

import (
	"context"

	"github.com/digisata/invitation-service/internal/entity"
)

type (
	InvitationRepository interface {
		Create(ctx context.Context, req entity.CreateInvitationRequest) error
		Update(ctx context.Context, req entity.UpdateInvitationRequest) error
		GetAll(ctx context.Context, req entity.GetAllInvitationRequest) ([]entity.Invitation, entity.Paging, error)
		GetByID(ctx context.Context, id string) (entity.Invitation, error)
		Delete(ctx context.Context, id string) error
	}

	InvitationLabelRepository interface {
		Create(ctx context.Context, req entity.CreateInvitationLabelRequest) error
		Update(ctx context.Context, req entity.UpdateInvitationLabelRequest) error
		GetAll(ctx context.Context, req entity.GetAllInvitationLabelRequest) ([]entity.InvitationLabel, entity.Paging, error)
		GetByID(ctx context.Context, id string) (entity.InvitationLabel, error)
		Delete(ctx context.Context, id string) error
	}

	InvitationCategoryRepository interface {
		Create(ctx context.Context, req entity.CreateInvitationCategoryRequest) error
		Update(ctx context.Context, req entity.UpdateInvitationCategoryRequest) error
		GetAll(ctx context.Context, req entity.GetAllInvitationCategoryRequest) ([]entity.InvitationCategory, entity.Paging, error)
		GetByID(ctx context.Context, id string) (entity.InvitationCategory, error)
		Delete(ctx context.Context, id string) error
	}
)
