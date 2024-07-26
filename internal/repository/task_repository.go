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

type TaskRepository struct {
	*postgres.Postgres
}

func NewTask(db *postgres.Postgres) *TaskRepository {
	return &TaskRepository{db}
}

func (r TaskRepository) Create(ctx context.Context, req entity.CreateTaskRequest) error {
	now := time.Now().UTC()
	sql, args, err := r.Builder.
		Insert("tasks").
		Columns("title, activity_id, is_active, priority, created_at, updated_at").
		Values(req.Title, req.ActivityID, req.IsActive, req.Priority, now, now).
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

func (r TaskRepository) Update(ctx context.Context, req entity.UpdateTaskRequest) error {
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
		Update("tasks").
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

func (r TaskRepository) GetAll(ctx context.Context, req entity.GetAllTaskRequest) ([]entity.Task, entity.Paging, error) {
	var (
		data   []entity.Task
		paging entity.Paging
	)

	baseQuery := r.Builder.
		Select("id, title, activity_id, is_active, priority, order_position, created_at, updated_at").
		From("tasks").
		Where(squirrel.Eq{"activity_id": req.ActivityID}).
		Where(squirrel.Eq{"deleted_at": nil})

	// Clone the base query for counting total rows
	countQuery := r.Builder.
		Select("COUNT(*)").
		From("tasks").
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

	if req.IsActive != nil {
		isFilterApplied = true
		baseQuery = baseQuery.Where(squirrel.Eq{"is_active": *req.IsActive})
		countQuery = countQuery.Where(squirrel.Eq{"is_active": *req.IsActive})
	}

	if req.Priority != nil {
		isFilterApplied = true
		baseQuery = baseQuery.Where(squirrel.Eq{"priority": *req.Priority})
		countQuery = countQuery.Where(squirrel.Eq{"priority": *req.Priority})
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
		baseQuery = baseQuery.OrderBy("title ASC")
	}

	if req.IsDescending != nil && *req.IsDescending {
		isFilterApplied = true
		baseQuery = baseQuery.OrderBy("title DESC")
	}

	if !isFilterApplied {
		baseQuery = baseQuery.OrderBy("order_position ASC")
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
		var task entity.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.ActivityID,
			&task.IsActive,
			&task.Priority,
			&task.Order,
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

func (r TaskRepository) GetByID(ctx context.Context, id string) (entity.Task, error) {
	var data entity.Task

	sql, args, err := r.Builder.
		Select("id, title, activity_id, is_active, priority, order_position, created_at, updated_at").
		From("tasks").
		Where(squirrel.Eq{"id": id}).
		Where(squirrel.Eq{"deleted_at": nil}).
		ToSql()
	if err != nil {
		return data, err
	}

	rows := r.Db.QueryRowContext(ctx, sql, args...)
	err = rows.Scan(
		&data.ID,
		&data.Title,
		&data.ActivityID,
		&data.IsActive,
		&data.Priority,
		&data.Order,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r TaskRepository) Delete(ctx context.Context, id string) error {
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
		Update("tasks").
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
