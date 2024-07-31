package entity

import "time"

type (
	Activity struct {
		ID        string
		Title     string
		Type      string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}

	CreateActivityRequest struct {
		Title string
		Type  string
	}

	UpdateActivityRequest struct {
		ID    string
		Title string `db:"title"`
		Type  string `db:"type"`
	}

	GetAllActivityRequest struct {
		Search *string
		Page   *int32
		Limit  *int32
	}
)
