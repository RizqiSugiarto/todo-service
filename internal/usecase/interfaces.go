package usecase

import (
	"context"

	"github.com/digisata/todo-service/internal/entity"
)

type (
	TaskRepository interface {
		Create(ctx context.Context, req entity.CreateTaskRequest) error
		Update(ctx context.Context, req entity.UpdateTaskRequest) error
		GetAll(ctx context.Context, req entity.GetAllTaskRequest) ([]entity.Task, entity.Paging, error)
		GetByID(ctx context.Context, id string) (entity.Task, error)
		Delete(ctx context.Context, id string) error
	}

	ActivityRepository interface {
		Create(ctx context.Context, req entity.CreateActivityRequest) error
		Update(ctx context.Context, req entity.UpdateActivityRequest) error
		GetAll(ctx context.Context, req entity.GetAllActivityRequest) ([]entity.Activity, entity.Paging, error)
		GetByID(ctx context.Context, id string) (entity.Activity, error)
		Delete(ctx context.Context, id string) error
	}

	TextRepository interface {
		Create(ctx context.Context, req entity.CreateTextRequest) error
		Update(ctx context.Context, req entity.UpdateTextRequest) error
		GetAll(ctx context.Context, req entity.GetAllTextRequest) ([]entity.Text, entity.Paging, error)
		GetByID(ctx context.Context, id string) (entity.Text, error)
		Delete(ctx context.Context, id string) error
	}
)
