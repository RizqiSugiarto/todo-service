package usecase

import (
	"context"

	"github.com/digisata/todo-service/internal/entity"
	"github.com/digisata/todo-service/internal/shared"
)

type TaskUseCase struct {
	taskRepository TaskRepository
}

func NewTask(taskRepository TaskRepository) *TaskUseCase {
	return &TaskUseCase{taskRepository: taskRepository}
}

func (u TaskUseCase) CreateTask(ctx context.Context, req entity.CreateTaskRequest) error {
	err := u.taskRepository.Create(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u TaskUseCase) UpdateTask(ctx context.Context, req entity.UpdateTaskRequest) error {
	err := u.taskRepository.Update(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u TaskUseCase) BatchUpdateTask(ctx context.Context, req []entity.UpdateTaskRequest) error {
	for _, task := range req {
		err := u.taskRepository.Update(ctx, task)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u TaskUseCase) GetTask(ctx context.Context, id string) (entity.Task, error) {
	var res entity.Task
	res, err := u.taskRepository.GetByID(ctx, id)
	if err != nil {
		return res, err
	}

	res.CreatedAt = shared.ConvertToJakartaTime(res.CreatedAt)
	res.UpdatedAt = shared.ConvertToJakartaTime(res.UpdatedAt)

	return res, nil
}

func (u TaskUseCase) GetAllTaskByActivityID(ctx context.Context, req entity.GetAllTaskRequest) ([]entity.Task, entity.Paging, error) {
	var (
		res    []entity.Task
		paging entity.Paging
	)
	res, paging, err := u.taskRepository.GetAll(ctx, req)
	if err != nil {
		return res, paging, err
	}

	for i := 0; i < len(res); i++ {
		res[i].CreatedAt = shared.ConvertToJakartaTime(res[i].CreatedAt)
		res[i].UpdatedAt = shared.ConvertToJakartaTime(res[i].UpdatedAt)
	}

	return res, paging, nil
}

func (u TaskUseCase) DeleteTask(ctx context.Context, id string) error {
	err := u.taskRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
