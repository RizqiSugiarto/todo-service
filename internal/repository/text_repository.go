package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/digisata/todo-service/internal/entity"
	"github.com/digisata/todo-service/internal/shared"
	"github.com/digisata/todo-service/pkg/postgres"
)

type TextRepository struct {
	*postgres.Postgres
}

func NewText(db *postgres.Postgres) *TextRepository {
	return &TextRepository{db}
}

func (r TextRepository) Create(ctx context.Context, req entity.CreateTextRequest) error {
	now := time.Now().UTC()
	sql, args, err := r.Builder.
		Insert("texts").
		Columns("text, activity_id, created_at, updated_at").
		Values(req.Text, req.ActivityID, now, now).
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

func (r TextRepository) Update(ctx context.Context, req entity.UpdateTextRequest) error {
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
		Update("texts").
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

func (r TextRepository) GetAll(ctx context.Context, req entity.GetAllTextRequest) ([]entity.Text, entity.Paging, error) {
	var (
		data   []entity.Text
		paging entity.Paging
	)

	baseQuery := r.Builder.
		Select("id, text, activity_id, created_at, updated_at").
		From("texts").
		Where(squirrel.Eq{"activity_id": req.ActivityID}).
		Where(squirrel.Eq{"deleted_at": nil})

	// Clone the base query for counting total rows
	countQuery := r.Builder.
		Select("COUNT(*)").
		From("texts").
		Where(squirrel.Eq{"activity_id": req.ActivityID}).
		Where(squirrel.Eq{"deleted_at": nil})

	var isFilterApplied bool

	// Apply search filter if present
	if req.Search != nil {
		isFilterApplied = true
		searchPattern := fmt.Sprintf("%%%s%%", *req.Search)
		baseQuery = baseQuery.Where(squirrel.ILike{"title": searchPattern})
		countQuery = countQuery.Where(squirrel.ILike{"title": searchPattern})
	}

	if req.IsNewest != nil && *req.IsNewest {
		isFilterApplied = true
		baseQuery = baseQuery.OrderBy("created_at DESC")
	}

	if req.IsOldest != nil && *req.IsOldest {
		isFilterApplied = true
		baseQuery = baseQuery.OrderBy("created_at ASC")
	}

	if req.IsAscending != nil && *req.IsAscending {
		isFilterApplied = true
		baseQuery = baseQuery.OrderBy("text ASC")
	}

	if req.IsDescending != nil && *req.IsDescending {
		isFilterApplied = true
		baseQuery = baseQuery.OrderBy("text DESC")
	}

	if !isFilterApplied {
		baseQuery = baseQuery.OrderBy("created_at ASC")
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
		var task entity.Text
		err := rows.Scan(
			&task.ID,
			&task.Text,
			&task.ActivityID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return data, paging, err
		}

		data = append(data, task)
	}

	return data, paging, nil
}

func (r TextRepository) GetByID(ctx context.Context, id string) (entity.Text, error) {
	var data entity.Text

	sql, args, err := r.Builder.
		Select("id, text, activity_id, created_at, updated_at").
		From("texts").
		Where(squirrel.Eq{"id": id}).
		Where(squirrel.Eq{"deleted_at": nil}).
		ToSql()
	if err != nil {
		return data, err
	}

	rows := r.Db.QueryRowContext(ctx, sql, args...)
	err = rows.Scan(
		&data.ID,
		&data.ActivityID,
		&data.Text,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r TextRepository) Delete(ctx context.Context, id string) error {
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
		Update("texts").
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
