package usecase

import (
	"context"

	"github.com/digisata/invitation-service/internal/entity"
	"github.com/digisata/invitation-service/internal/shared"
)

type InvitationLabelUseCase struct {
	invitationCategoryRepository InvitationLabelRepository
}

func NewInvitationLabel(invitationCategoryRepository InvitationLabelRepository) *InvitationLabelUseCase {
	return &InvitationLabelUseCase{invitationCategoryRepository: invitationCategoryRepository}
}

func (u InvitationLabelUseCase) CreateInvitationLabel(ctx context.Context, req entity.CreateInvitationLabelRequest) error {
	err := u.invitationCategoryRepository.Create(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u InvitationLabelUseCase) UpdateInvitationLabel(ctx context.Context, req entity.UpdateInvitationLabelRequest) error {
	err := u.invitationCategoryRepository.Update(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u InvitationLabelUseCase) GetInvitationLabel(ctx context.Context, id string) (entity.InvitationLabel, error) {
	var res entity.InvitationLabel
	res, err := u.invitationCategoryRepository.GetByID(ctx, id)
	if err != nil {
		return res, err
	}

	res.CreatedAt = shared.ConvertToJakartaTime(res.CreatedAt)
	res.UpdatedAt = shared.ConvertToJakartaTime(res.UpdatedAt)

	return res, nil
}

func (u InvitationLabelUseCase) GetAllInvitationLabelByUserID(ctx context.Context, req entity.GetAllInvitationLabelRequest) ([]entity.InvitationLabel, entity.Paging, error) {
	var (
		res    []entity.InvitationLabel
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

func (u InvitationLabelUseCase) DeleteInvitationLabel(ctx context.Context, id string) error {
	err := u.invitationCategoryRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
