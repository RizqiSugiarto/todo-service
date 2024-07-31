package handler

import (
	"context"

	"github.com/digisata/todo-service/internal/entity"
)

type (
	TaskUseCase interface {
		CreateTask(ctx context.Context, req entity.CreateTaskRequest) error
		UpdateTask(ctx context.Context, req entity.UpdateTaskRequest) error
		BatchUpdateTask(ctx context.Context, req []entity.UpdateTaskRequest) error
		GetTask(ctx context.Context, id string) (entity.Task, error)
		GetAllTaskByActivityID(ctx context.Context, req entity.GetAllTaskRequest) ([]entity.Task, entity.Paging, error)
		DeleteTask(ctx context.Context, id string) error
	}

	ActivityUseCase interface {
		CreateActivity(ctx context.Context, req entity.CreateActivityRequest) error
		UpdateActivity(ctx context.Context, req entity.UpdateActivityRequest) error
		GetActivity(ctx context.Context, id string) (entity.Activity, error)
		GetAllActivity(ctx context.Context, req entity.GetAllActivityRequest) ([]entity.Activity, entity.Paging, error)
		DeleteActivity(ctx context.Context, id string) error
	}

	TextUseCase interface {
		CreateText(ctx context.Context, req entity.CreateTextRequest) error
		UpdateText(ctx context.Context, req entity.UpdateTextRequest) error
		GetText(ctx context.Context, id string) (entity.Text, error)
		GetAllTextByActivityID(ctx context.Context, req entity.GetAllTextRequest) ([]entity.Text, entity.Paging, error)
		DeleteText(ctx context.Context, id string) error
	}
)
