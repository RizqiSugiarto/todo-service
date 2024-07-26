package usecase

import (
	"context"

	"github.com/digisata/todo-service/internal/entity"
	"github.com/digisata/todo-service/internal/shared"
)

type ActivityUseCase struct {
	activityRepository ActivityRepository
}

func NewActivity(activityRepository ActivityRepository) *ActivityUseCase {
	return &ActivityUseCase{activityRepository: activityRepository}
}

func (u ActivityUseCase) CreateActivity(ctx context.Context, req entity.CreateActivityRequest) error {
	err := u.activityRepository.Create(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u ActivityUseCase) UpdateActivity(ctx context.Context, req entity.UpdateActivityRequest) error {
	err := u.activityRepository.Update(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u ActivityUseCase) GetActivity(ctx context.Context, id string) (entity.Activity, error) {
	var res entity.Activity
	res, err := u.activityRepository.GetByID(ctx, id)
	if err != nil {
		return res, err
	}

	res.CreatedAt = shared.ConvertToJakartaTime(res.CreatedAt)
	res.UpdatedAt = shared.ConvertToJakartaTime(res.UpdatedAt)

	return res, nil
}

func (u ActivityUseCase) GetAllActivity(ctx context.Context, req entity.GetAllActivityRequest) ([]entity.Activity, entity.Paging, error) {
	var (
		res    []entity.Activity
		paging entity.Paging
	)
	res, paging, err := u.activityRepository.GetAll(ctx, req)
	if err != nil {
		return res, paging, err
	}

	for i := 0; i < len(res); i++ {
		res[i].CreatedAt = shared.ConvertToJakartaTime(res[i].CreatedAt)
		res[i].UpdatedAt = shared.ConvertToJakartaTime(res[i].UpdatedAt)
	}

	return res, paging, nil
}

func (u ActivityUseCase) DeleteActivity(ctx context.Context, id string) error {
	err := u.activityRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
