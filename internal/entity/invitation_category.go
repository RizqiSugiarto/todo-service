package entity

import "time"

type (
	InvitationCategory struct {
		ID        string
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}

	CreateInvitationCategoryRequest struct {
		Name string
	}

	UpdateInvitationCategoryRequest struct {
		ID   string
		Name string `db:"name"`
	}

	GetAllInvitationCategoryRequest struct {
		Search *string
		Page   *int32
		Limit  *int32
	}
)
