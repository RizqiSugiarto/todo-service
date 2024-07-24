package entity

import "time"

type (
	InvitationLabel struct {
		ID        string
		UserID    string
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}

	CreateInvitationLabelRequest struct {
		UserID string
		Name   string
	}

	UpdateInvitationLabelRequest struct {
		ID   string
		Name string `db:"name"`
	}

	GetAllInvitationLabelRequest struct {
		UserID string
		Search *string
		Page   *int32
		Limit  *int32
	}
)
