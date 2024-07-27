package usecase

import (
	"context"

	"github.com/digisata/todo-service/internal/entity"
	"github.com/digisata/todo-service/internal/shared"
)

type TextUseCase struct {
	textRepository TextRepository
}

func NewText(textRepository TextRepository) *TextUseCase {
	return &TextUseCase{textRepository: textRepository}
}

func (u TextUseCase) CreateText(ctx context.Context, req entity.CreateTextRequest) error {
	err := u.textRepository.Create(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u TextUseCase) UpdateText(ctx context.Context, req entity.UpdateTextRequest) error {
	err := u.textRepository.Update(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

// func (u TextUseCase) BatchUpdateText(ctx context.Context, req []entity.UpdateTextRequest) error {
// 	for _, task := range req {
// 		err := u.textRepository.Update(ctx, task)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (u TextUseCase) GetText(ctx context.Context, id string) (entity.Text, error) {
	var res entity.Text
	res, err := u.textRepository.GetByID(ctx, id)
	if err != nil {
		return res, err
	}

	res.CreatedAt = shared.ConvertToJakartaTime(res.CreatedAt)
	res.UpdatedAt = shared.ConvertToJakartaTime(res.UpdatedAt)

	return res, nil
}

func (u TextUseCase) GetAllTextByActivityID(ctx context.Context, req entity.GetAllTextRequest) ([]entity.Text, entity.Paging, error) {
	var (
		res    []entity.Text
		paging entity.Paging
	)
	res, paging, err := u.textRepository.GetAll(ctx, req)
	if err != nil {
		return res, paging, err
	}

	for i := 0; i < len(res); i++ {
		res[i].CreatedAt = shared.ConvertToJakartaTime(res[i].CreatedAt)
		res[i].UpdatedAt = shared.ConvertToJakartaTime(res[i].UpdatedAt)
	}

	return res, paging, nil
}

func (u TextUseCase) DeleteText(ctx context.Context, id string) error {
	err := u.textRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
