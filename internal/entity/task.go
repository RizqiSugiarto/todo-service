package entity

import "time"

type (
	Task struct {
		ID         string
		Title      string
		ActivityID string
		IsActive   bool
		Priority   int
		Order      int
		CreatedAt  time.Time
		UpdatedAt  time.Time
		DeletedAt  *time.Time
	}

	CreateTaskRequest struct {
		Title      string
		ActivityID string
		IsActive   *bool
		Priority   int
	}

	UpdateTaskRequest struct {
		ID       string
		Title    *string `db:"title"`
		IsActive *bool   `db:"is_active"`
		Priority *int    `db:"priority"`
		Order    *int    `db:"order_position"`
	}

	GetAllTaskRequest struct {
		ActivityID   string
		IsActive     *bool
		Priority     *int
		IsNewest     *bool
		IsOldest     *bool
		IsAscending  *bool
		IsDescending *bool
		Search       *string
		Page         *int32
		Limit        *int32
	}
)
