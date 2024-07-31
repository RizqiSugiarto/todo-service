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

type ActivityRepository struct {
	*postgres.Postgres
}

func NewActivity(db *postgres.Postgres) *ActivityRepository {
	return &ActivityRepository{db}
}

func (r ActivityRepository) Create(ctx context.Context, req entity.CreateActivityRequest) error {
	var activityId string

	now := time.Now().UTC()
	sql, args, err := r.Builder.
		Insert("activities").
		Columns("title, type, created_at, updated_at").
		Values(req.Title, req.Type, now, now).Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	err = r.Db.QueryRowContext(ctx, sql, args...).Scan(&activityId)
	if err != nil {
		return err
	}

	if req.Type == "activity_text" {
		sql, args, err := r.Builder.
			Insert("texts").
			Columns("text, activity_id, created_at, updated_at").
			Values(`<p class="default-text">Fill your note ....</p>`, activityId, now, now).
			ToSql()
		if err != nil {
			return err
		}

		_, err = r.Db.ExecContext(ctx, sql, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r ActivityRepository) Update(ctx context.Context, req entity.UpdateActivityRequest) error {
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
		Update("activities").
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

func (r ActivityRepository) GetAll(ctx context.Context, req entity.GetAllActivityRequest) ([]entity.Activity, entity.Paging, error) {
	var (
		data   []entity.Activity
		paging entity.Paging
	)

	baseQuery := r.Builder.
		Select("id, title, type, created_at, updated_at").
		From("activities").
		Where(squirrel.Eq{"deleted_at": nil})

	// Clone the base query for counting total rows
	countQuery := r.Builder.
		Select("COUNT(*)").
		From("activities a").
		Where(squirrel.Eq{"a.deleted_at": nil})

	// Apply search filter if present
	if req.Search != nil {
		searchPattern := fmt.Sprintf("%%%s%%", *req.Search)
		baseQuery = baseQuery.Where(squirrel.ILike{"a.title": searchPattern})
		countQuery = countQuery.Where(squirrel.ILike{"a.title": searchPattern})
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
		var activity entity.Activity
		if err := rows.Scan(&activity.ID, &activity.Title, &activity.Type, &activity.CreatedAt, &activity.UpdatedAt); err != nil {
			return data, paging, err
		}
		data = append(data, activity)
	}

	return data, paging, nil
}

func (r ActivityRepository) GetByID(ctx context.Context, id string) (entity.Activity, error) {
	var data entity.Activity

	sql, args, err := r.Builder.
		Select("id, title, type,created_at, updated_at").
		From("activities").
		Where(squirrel.Eq{"id": id}).
		Where(squirrel.Eq{"deleted_at": nil}).
		ToSql()
	if err != nil {
		return data, err
	}

	row := r.Db.QueryRowContext(ctx, sql, args...)
	err = row.Scan(&data.ID, &data.Title, &data.Type, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r ActivityRepository) Delete(ctx context.Context, id string) error {
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
		Update("activities").
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
