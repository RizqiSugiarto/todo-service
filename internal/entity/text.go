package entity

import "time"

type (
	Text struct {
		ID         string
		ActivityID string
		Text       string
		CreatedAt  time.Time
		UpdatedAt  time.Time
		DeletedAt  *time.Time
	}

	CreateTextRequest struct {
		ActivityID string
		Text       string
	}

	UpdateTextRequest struct {
		ID   string
		Text *string `db:"text"`
	}

	GetAllTextRequest struct {
		ActivityID   string
		IsNewest     *bool
		IsOldest     *bool
		IsAscending  *bool
		IsDescending *bool
		Search       *string
		Page         *int32
		Limit        *int32
	}
)
