package entity

import "time"

type (
	Invitation struct {
		ID                          string
		UserID                      string
		Name                        string
		InvitationLabelID           *string
		InvitationLabelUserID       *string
		InvitationLabelName         *string
		InvitationLabelCreatedAt    *time.Time
		InvitationLabelUpdatedAt    *time.Time
		InvitationLabelDeletedAt    *time.Time
		InvitationCategoryID        *string
		InvitationCategoryName      *string
		InvitationCategoryCreatedAt *time.Time
		InvitationCategoryUpdatedAt *time.Time
		InvitationCategoryDeletedAt *time.Time
		IsOpen                      bool
		OpenAt                      *time.Time
		IsComing                    *bool
		IsSendMoney                 bool
		IsSendGift                  bool
		IsCheckIn                   bool
		CheckInAt                   *time.Time
		CreatedAt                   time.Time
		UpdatedAt                   time.Time
		DeletedAt                   *time.Time
	}

	CreateInvitationRequest struct {
		UserID               string
		Name                 string
		InvitationLabelID    *string
		InvitationCategoryID *string
	}

	UpdateInvitationRequest struct {
		ID                   string
		Name                 *string    `db:"name"`
		InvitationLabelID    *string    `db:"invitation_label_id"`
		InvitationCategoryID *string    `db:"invitation_category_id"`
		IsOpen               *bool      `db:"is_open"`
		OpenAt               *time.Time `db:"open_at"`
		IsComing             *bool      `db:"is_coming"`
		IsSendMoney          *bool      `db:"is_send_money"`
		IsSendGift           *bool      `db:"is_send_gift"`
		IsCheckIn            *bool      `db:"is_check_in"`
		CheckInAt            *time.Time `db:"check_in_at"`
	}

	GetAllInvitationRequest struct {
		UserID               string
		Search               *string
		Page                 *int32
		Limit                *int32
		IsOpen               *bool
		IsComing             *bool
		IsSendMoney          *bool
		IsSendGift           *bool
		IsCheckIn            *bool
		InvitationLabels     []string
		InvitationCategories []string
	}

	InvitationEventRequest struct {
		ID        string     `json:"id"`
		IsComing  *bool      `json:"is_coming"`
		OpenAt    *time.Time `json:"open_at"`
		CheckInAt *time.Time `json:"check_in_at"`
	}

	VerifyInvitationResponse struct {
		UserID       string `json:"user_id"`
		InvitationID string `json:"invitation_id"`
	}
)
