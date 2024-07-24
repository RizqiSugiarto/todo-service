package usecase

import (
	"context"

	"github.com/digisata/invitation-service/internal/entity"
	"github.com/digisata/invitation-service/internal/shared"
)

type InvitationCategoryUseCase struct {
	invitationCategoryRepository InvitationCategoryRepository
}

func NewInvitationCategory(invitationCategoryRepository InvitationCategoryRepository) *InvitationCategoryUseCase {
	return &InvitationCategoryUseCase{invitationCategoryRepository: invitationCategoryRepository}
}

func (u InvitationCategoryUseCase) CreateInvitationCategory(ctx context.Context, req entity.CreateInvitationCategoryRequest) error {
	err := u.invitationCategoryRepository.Create(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u InvitationCategoryUseCase) UpdateInvitationCategory(ctx context.Context, req entity.UpdateInvitationCategoryRequest) error {
	err := u.invitationCategoryRepository.Update(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u InvitationCategoryUseCase) GetInvitationCategory(ctx context.Context, id string) (entity.InvitationCategory, error) {
	var res entity.InvitationCategory
	res, err := u.invitationCategoryRepository.GetByID(ctx, id)
	if err != nil {
		return res, err
	}

	res.CreatedAt = shared.ConvertToJakartaTime(res.CreatedAt)
	res.UpdatedAt = shared.ConvertToJakartaTime(res.UpdatedAt)

	return res, nil
}

func (u InvitationCategoryUseCase) GetAllInvitationCategory(ctx context.Context, req entity.GetAllInvitationCategoryRequest) ([]entity.InvitationCategory, entity.Paging, error) {
	var (
		res    []entity.InvitationCategory
		paging entity.Paging
	)
	res, paging, err := u.invitationCategoryRepository.GetAll(ctx, req)
	if err != nil {
		return res, paging, err
	}

	for i := 0; i < len(res); i++ {
		res[i].CreatedAt = shared.ConvertToJakartaTime(res[i].CreatedAt)
		res[i].UpdatedAt = shared.ConvertToJakartaTime(res[i].UpdatedAt)
	}

	return res, paging, nil
}

func (u InvitationCategoryUseCase) DeleteInvitationCategory(ctx context.Context, id string) error {
	err := u.invitationCategoryRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
