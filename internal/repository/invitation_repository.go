package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/digisata/invitation-service/internal/entity"
	"github.com/digisata/invitation-service/internal/shared"
	"github.com/digisata/invitation-service/pkg/postgres"
)

type InvitationRepository struct {
	*postgres.Postgres
}

func NewInvitation(db *postgres.Postgres) *InvitationRepository {
	return &InvitationRepository{db}
}

func (r InvitationRepository) Create(ctx context.Context, req entity.CreateInvitationRequest) error {
	now := time.Now().UTC()
	sql, args, err := r.Builder.
		Insert("invitations").
		Columns("user_id, name, invitation_label_id, invitation_category_id, created_at, updated_at").
		Values(req.UserID, req.Name, req.InvitationLabelID, req.InvitationCategoryID, now, now).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.Db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r InvitationRepository) Update(ctx context.Context, req entity.UpdateInvitationRequest) error {
	tx, err := r.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	updateValue := shared.CreateUpdateValueMap(req)

	sql, args, err := r.Builder.
		Update("invitations").
		SetMap(updateValue).
		Where(squirrel.Eq{"id": req.ID}).
		Where(squirrel.Eq{"deleted_at": nil}).
		ToSql()
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("data not found")
	}

	return nil
}

func (r InvitationRepository) GetAll(ctx context.Context, req entity.GetAllInvitationRequest) ([]entity.Invitation, entity.Paging, error) {
	var (
		data   []entity.Invitation
		paging entity.Paging
	)

	baseQuery := r.Builder.
		Select("i.id, i.user_id, i.name, il.id, il.user_id, il.name, il.created_at, il.updated_at, ic.id, ic.name, ic.created_at, ic.updated_at, i.is_open, i.open_at, i.is_coming, i.is_send_money, i.is_send_gift, i.is_check_in, i.check_in_at, i.created_at, i.updated_at").
		From("invitations i").
		LeftJoin("invitation_labels il on il.id = i.invitation_label_id").
		LeftJoin("invitation_categories ic on ic.id = i.invitation_category_id").
		Where(squirrel.Eq{"i.user_id": req.UserID}).
		Where(squirrel.Eq{"i.deleted_at": nil})

	// Clone the base query for counting total rows
	countQuery := r.Builder.
		Select("COUNT(*)").
		From("invitations").
		Where(squirrel.Eq{"user_id": req.UserID}).
		Where(squirrel.Eq{"deleted_at": nil})

	// Apply search filter if present
	if req.Search != nil {
		searchPattern := fmt.Sprintf("%%%s%%", *req.Search)
		baseQuery = baseQuery.Where(squirrel.ILike{"i.name": searchPattern})
		countQuery = countQuery.Where(squirrel.ILike{"name": searchPattern})
	}

	if req.IsOpen != nil {
		baseQuery = baseQuery.Where(squirrel.Eq{"i.is_open": *req.IsOpen})
		countQuery = countQuery.Where(squirrel.Eq{"is_open": *req.IsOpen})
	}

	if req.IsComing != nil {
		baseQuery = baseQuery.Where(squirrel.Eq{"i.is_coming": *req.IsComing})
		countQuery = countQuery.Where(squirrel.Eq{"is_coming": *req.IsComing})
	}

	if req.IsSendMoney != nil {
		baseQuery = baseQuery.Where(squirrel.Eq{"i.is_send_money": *req.IsSendMoney})
		countQuery = countQuery.Where(squirrel.Eq{"is_send_money": *req.IsSendMoney})
	}

	if req.IsSendGift != nil {
		baseQuery = baseQuery.Where(squirrel.Eq{"i.is_send_gift": *req.IsSendGift})
		countQuery = countQuery.Where(squirrel.Eq{"is_send_gift": *req.IsSendGift})
	}

	if req.IsCheckIn != nil {
		baseQuery = baseQuery.Where(squirrel.Eq{"i.is_check_in": *req.IsCheckIn})
		countQuery = countQuery.Where(squirrel.Eq{"is_check_in": *req.IsCheckIn})
	}

	// Apply InvitationLabels filter if present
	if req.InvitationLabels != nil && len(req.InvitationLabels) > 0 {
		baseQuery = baseQuery.Where(squirrel.Eq{"i.invitation_label_id": req.InvitationLabels})
		countQuery = countQuery.Where(squirrel.Eq{"invitation_label_id": req.InvitationLabels})
	}

	// Apply InvitationCategories filter if present
	if req.InvitationCategories != nil && len(req.InvitationCategories) > 0 {
		baseQuery = baseQuery.Where(squirrel.Eq{"i.invitation_category_id": req.InvitationCategories})
		countQuery = countQuery.Where(squirrel.Eq{"invitation_category_id": req.InvitationCategories})
	}

	// Get the total count of rows that match the query
	totalRowsSql, totalRowsArgs, err := countQuery.ToSql()
	if err != nil {
		return data, paging, err
	}

	var totalRows int32
	err = r.Db.QueryRowContext(ctx, totalRowsSql, totalRowsArgs...).Scan(&totalRows)
	if err != nil {
		return data, paging, err
	}

	// Calculate total pages
	if req.Limit != nil && *req.Limit > 0 {
		paging.TotalPage = (totalRows + *req.Limit - 1) / *req.Limit
	} else {
		paging.TotalPage = 1
	}

	// Set current page
	if req.Page != nil && *req.Page > 0 {
		paging.CurrentPage = *req.Page
	} else {
		paging.CurrentPage = 1
	}

	paging.Count = int32(totalRows)

	// Apply pagination if both page and limit are provided
	if req.Page != nil && req.Limit != nil && *req.Limit > 0 {
		offset := (*req.Page - 1) * *req.Limit
		baseQuery = baseQuery.Limit(uint64(*req.Limit)).Offset(uint64(offset))
	}

	// Execute the query to get paginated data
	sql, args, err := baseQuery.ToSql()
	if err != nil {
		return data, paging, err
	}

	rows, err := r.Db.QueryContext(ctx, sql, args...)
	if err != nil {
		return data, paging, err
	}
	defer rows.Close()

	for rows.Next() {
		var invitation entity.Invitation
		err := rows.Scan(
			&invitation.ID,
			&invitation.UserID,
			&invitation.Name,
			&invitation.InvitationLabelID,
			&invitation.InvitationLabelUserID,
			&invitation.InvitationLabelName,
			&invitation.InvitationLabelCreatedAt,
			&invitation.InvitationLabelUpdatedAt,
			&invitation.InvitationCategoryID,
			&invitation.InvitationCategoryName,
			&invitation.InvitationCategoryCreatedAt,
			&invitation.InvitationCategoryUpdatedAt,
			&invitation.IsOpen,
			&invitation.OpenAt,
			&invitation.IsComing,
			&invitation.IsSendMoney,
			&invitation.IsSendGift,
			&invitation.IsCheckIn,
			&invitation.CheckInAt,
			&invitation.CreatedAt,
			&invitation.UpdatedAt,
		)
		if err != nil {
			return data, paging, err
		}

		data = append(data, invitation)
	}

	return data, paging, nil
}

func (r InvitationRepository) GetByID(ctx context.Context, id string) (entity.Invitation, error) {
	var data entity.Invitation

	sql, args, err := r.Builder.
		Select("i.id, i.user_id, i.name, il.id, il.user_id, il.name, il.created_at, il.updated_at, ic.id, ic.name, ic.created_at, ic.updated_at, i.is_open, i.open_at, i.is_coming, i.is_send_money, i.is_send_gift, i.is_check_in, i.check_in_at, i.created_at, i.updated_at").
		From("invitations i").
		LeftJoin("invitation_labels il on il.id = i.invitation_label_id").
		LeftJoin("invitation_categories ic on ic.id = i.invitation_category_id").
		Where(squirrel.Eq{"i.id": id}).
		Where(squirrel.Eq{"i.deleted_at": nil}).
		ToSql()
	if err != nil {
		return data, err
	}

	rows := r.Db.QueryRowContext(ctx, sql, args...)
	err = rows.Scan(
		&data.ID,
		&data.UserID,
		&data.Name,
		&data.InvitationLabelID,
		&data.InvitationLabelUserID,
		&data.InvitationLabelName,
		&data.InvitationLabelCreatedAt,
		&data.InvitationLabelUpdatedAt,
		&data.InvitationCategoryID,
		&data.InvitationCategoryName,
		&data.InvitationCategoryCreatedAt,
		&data.InvitationCategoryUpdatedAt,
		&data.IsOpen,
		&data.OpenAt,
		&data.IsComing,
		&data.IsSendMoney,
		&data.IsSendGift,
		&data.IsCheckIn,
		&data.CheckInAt,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r InvitationRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	deleteValue := map[string]interface{}{
		"deleted_at": time.Now().UTC(),
	}

	sql, args, err := r.Builder.
		Update("invitations").
		SetMap(deleteValue).
		Where(squirrel.Eq{"id": id}).
		Where(squirrel.Eq{"deleted_at": nil}).
		ToSql()
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("data not found")
	}

	return nil
}
