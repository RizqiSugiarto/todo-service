package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/digisata/invitation-service/config"
	"github.com/digisata/invitation-service/internal/entity"
	"github.com/digisata/invitation-service/internal/helper"
	"github.com/digisata/invitation-service/internal/shared"
)

type InvitationUseCase struct {
	invitationCategoryRepository InvitationRepository
	cfg                          *config.Config
}

func NewInvitation(cfg *config.Config, invitationCategoryRepository InvitationRepository) *InvitationUseCase {
	return &InvitationUseCase{
		cfg:                          cfg,
		invitationCategoryRepository: invitationCategoryRepository,
	}
}

func (u InvitationUseCase) CreateInvitation(ctx context.Context, req entity.CreateInvitationRequest) error {
	err := u.invitationCategoryRepository.Create(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u InvitationUseCase) UpdateInvitation(ctx context.Context, req entity.UpdateInvitationRequest) error {
	err := u.invitationCategoryRepository.Update(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u InvitationUseCase) GetInvitation(ctx context.Context, id string) (entity.Invitation, error) {
	var res entity.Invitation
	res, err := u.invitationCategoryRepository.GetByID(ctx, id)
	if err != nil {
		return res, err
	}

	if res.InvitationLabelID != nil {
		createdAt := shared.ConvertToJakartaTime(*res.InvitationLabelCreatedAt)
		res.InvitationLabelCreatedAt = &createdAt

		updatedAt := shared.ConvertToJakartaTime(*res.InvitationLabelUpdatedAt)
		res.InvitationLabelUpdatedAt = &updatedAt
	}

	if res.InvitationCategoryID != nil {
		createdAt := shared.ConvertToJakartaTime(*res.InvitationCategoryCreatedAt)
		res.InvitationCategoryCreatedAt = &createdAt

		updatedAt := shared.ConvertToJakartaTime(*res.InvitationCategoryUpdatedAt)
		res.InvitationCategoryUpdatedAt = &updatedAt
	}

	if res.OpenAt != nil {
		openAt := shared.ConvertToJakartaTime(*res.OpenAt)
		res.OpenAt = &openAt
	}

	if res.CheckInAt != nil {
		checkInAt := shared.ConvertToJakartaTime(*res.CheckInAt)
		res.CheckInAt = &checkInAt
	}

	res.CreatedAt = shared.ConvertToJakartaTime(res.CreatedAt)
	res.UpdatedAt = shared.ConvertToJakartaTime(res.UpdatedAt)

	return res, nil
}

func (u InvitationUseCase) GetAllInvitationByUserID(ctx context.Context, req entity.GetAllInvitationRequest) ([]entity.Invitation, entity.Paging, error) {
	var (
		res    []entity.Invitation
		paging entity.Paging
	)

	res, paging, err := u.invitationCategoryRepository.GetAll(ctx, req)
	if err != nil {
		return res, paging, err
	}

	for i := 0; i < len(res); i++ {
		if res[i].InvitationLabelID != nil {
			createdAt := shared.ConvertToJakartaTime(*res[i].InvitationLabelCreatedAt)
			res[i].InvitationLabelCreatedAt = &createdAt

			updatedAt := shared.ConvertToJakartaTime(*res[i].InvitationLabelUpdatedAt)
			res[i].InvitationLabelUpdatedAt = &updatedAt
		}

		if res[i].InvitationCategoryID != nil {
			createdAt := shared.ConvertToJakartaTime(*res[i].InvitationCategoryCreatedAt)
			res[i].InvitationCategoryCreatedAt = &createdAt

			updatedAt := shared.ConvertToJakartaTime(*res[i].InvitationCategoryUpdatedAt)
			res[i].InvitationCategoryUpdatedAt = &updatedAt
		}

		if res[i].OpenAt != nil {
			openAt := shared.ConvertToJakartaTime(*res[i].OpenAt)
			res[i].OpenAt = &openAt
		}

		if res[i].CheckInAt != nil {
			checkInAt := shared.ConvertToJakartaTime(*res[i].CheckInAt)
			res[i].CheckInAt = &checkInAt
		}

		res[i].CreatedAt = shared.ConvertToJakartaTime(res[i].CreatedAt)
		res[i].UpdatedAt = shared.ConvertToJakartaTime(res[i].UpdatedAt)
	}

	return res, paging, nil
}

func (u InvitationUseCase) DeleteInvitation(ctx context.Context, id string) error {
	err := u.invitationCategoryRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u InvitationUseCase) GetInvitationLink(ctx context.Context, id string) (string, error) {
	var link string

	res, err := u.invitationCategoryRepository.GetByID(ctx, id)
	if err != nil {
		return link, err
	}

	link, err = helper.GenerateInvitationLink(res, u.cfg.EncryptionKey, u.cfg.WeddingInvitationBaseUrl, "jardin-eldeberry")
	if err != nil {
		return link, err
	}

	return link, nil
}

func (u InvitationUseCase) VerifyInvitation(ctx context.Context, id string) (entity.VerifyInvitationResponse, error) {
	var res entity.VerifyInvitationResponse

	plainText, err := helper.Decrypt(u.cfg.EncryptionKey, id)
	if err != nil {
		return res, err
	}

	arr := strings.Split(plainText, "~")

	if len(arr) != 2 {
		return res, errors.New("malformed encryption")
	}

	invitationData, err := u.invitationCategoryRepository.GetByID(ctx, arr[1])
	if err != nil {
		return res, err
	}

	if invitationData.UserID != arr[0] {
		return res, errors.New("the invitation is not in the list")
	}

	res = entity.VerifyInvitationResponse{
		UserID:       arr[0],
		InvitationID: arr[1],
	}

	return res, nil
}
