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

type InvitationCategoryRepository struct {
	*postgres.Postgres
}

func NewInvitationCategory(db *postgres.Postgres) *InvitationCategoryRepository {
	return &InvitationCategoryRepository{db}
}

func (r InvitationCategoryRepository) Create(ctx context.Context, req entity.CreateInvitationCategoryRequest) error {
	now := time.Now().UTC()
	sql, args, err := r.Builder.
		Insert("invitation_categories").
		Columns("name, created_at, updated_at").
		Values(req.Name, now, now).
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

func (r InvitationCategoryRepository) Update(ctx context.Context, req entity.UpdateInvitationCategoryRequest) error {
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
		Update("invitation_categories").
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

func (r InvitationCategoryRepository) GetAll(ctx context.Context, req entity.GetAllInvitationCategoryRequest) ([]entity.InvitationCategory, entity.Paging, error) {
	var (
		data   []entity.InvitationCategory
		paging entity.Paging
	)

	baseQuery := r.Builder.
		Select("id, name, created_at, updated_at").
		From("invitation_categories").
		Where(squirrel.Eq{"deleted_at": nil})

	// Clone the base query for counting total rows
	countQuery := r.Builder.
		Select("COUNT(*)").
		From("invitation_categories").
		Where(squirrel.Eq{"deleted_at": nil})

	// Apply search filter if present
	if req.Search != nil {
		searchPattern := fmt.Sprintf("%%%s%%", *req.Search)
		baseQuery = baseQuery.Where(squirrel.ILike{"name": searchPattern})
		countQuery = countQuery.Where(squirrel.ILike{"name": searchPattern})
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
		var invitationCategory entity.InvitationCategory
		if err := rows.Scan(&invitationCategory.ID, &invitationCategory.Name, &invitationCategory.CreatedAt, &invitationCategory.UpdatedAt); err != nil {
			return data, paging, err
		}
		data = append(data, invitationCategory)
	}

	return data, paging, nil
}

func (r InvitationCategoryRepository) GetByID(ctx context.Context, id string) (entity.InvitationCategory, error) {
	var data entity.InvitationCategory

	sql, args, err := r.Builder.
		Select("id, name, created_at, updated_at").
		From("invitation_categories").
		Where(squirrel.Eq{"id": id}).
		Where(squirrel.Eq{"deleted_at": nil}).
		ToSql()
	if err != nil {
		return data, err
	}

	rows := r.Db.QueryRowContext(ctx, sql, args...)
	err = rows.Scan(&data.ID, &data.Name, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r InvitationCategoryRepository) Delete(ctx context.Context, id string) error {
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
		Update("invitation_categories").
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
