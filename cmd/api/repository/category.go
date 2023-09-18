package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/dto"

	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/domain"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/lib/helperdb"
)

type CategoryRepository interface {
	Create(ctx context.Context, tx *sql.Tx, category domain.Category) (int64, error)
	FindAll(ctx context.Context, tx helperdb.Tx, params dto.SearchParams) ([]domain.Category, int64, error)
	FindByID(ctx context.Context, tx helperdb.Tx, id int64) (domain.Category, error)
	FindByName(ctx context.Context, tx helperdb.Tx, name string) (domain.Category, error)
	Update(ctx context.Context, tx *sql.Tx, category domain.Category) (int64, error)
	Delete(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
}

type categoryRepository struct {
}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{}
}

const (
	createCategoryQuery     = "INSERT INTO categories (name, created_at) VALUES ( ?, NOW())"
	findAllCategoryQuery    = "SELECT id, name, created_at FROM categories"
	findByIDCategoryQuery   = "SELECT id, name, created_at FROM categories WHERE id = ?"
	findByNameCategoryQuery = "SELECT id, name, created_at FROM categories WHERE name LIKE ?"
	updateCategoryQuery     = "UPDATE categories SET name = ? WHERE id = ?"
	deleteCategoryQuery     = "DELETE FROM categories WHERE id = ?"
)

func (c *categoryRepository) Create(ctx context.Context, tx *sql.Tx, category domain.Category) (int64, error) {
	res, err := tx.ExecContext(
		ctx,
		createCategoryQuery,
		category.Name,
	)
	if err != nil {
		return 0, domain.NewInternalError(err.Error(), err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, domain.NewInternalError("fail to get last insert id", err)
	}
	return id, nil
}

func (c *categoryRepository) FindAll(ctx context.Context, tx helperdb.Tx, params dto.SearchParams) ([]domain.Category, int64, error) {
	query, queryParams := GetSearchQuery(params, findAllCategoryQuery)
	var categories []domain.Category
	var total int64
	err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM categories").Scan(&total)
	if err != nil {
		return nil, 0, domain.NewInternalError(err.Error(), err)
	}
	rows, err := tx.QueryContext(ctx, query.String(), queryParams...)
	if err != nil {
		return nil, 0, domain.NewInternalError("fail to get categories from db", err)
	}
	defer rows.Close()
	for rows.Next() {
		var category domain.Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, 0, domain.NewInternalError("fail to scan category", err)
		}
		categories = append(categories, category)
	}
	return categories, total, nil
}

func (c *categoryRepository) FindByID(ctx context.Context, tx helperdb.Tx, id int64) (domain.Category, error) {
	row := tx.QueryRow(findByIDCategoryQuery, id)
	var categories domain.Category
	err := row.Scan(
		&categories.ID,
		&categories.Name,
		&categories.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Category{}, domain.NewNotFoundError(fmt.Sprintf("category with id %d not found", id), err)
		}
		return domain.Category{}, domain.NewInternalError(err.Error(), err)
	}

	return categories, nil
}

func (c *categoryRepository) FindByName(ctx context.Context, tx helperdb.Tx, name string) (domain.Category, error) {
	row := tx.QueryRow(findByNameCategoryQuery, name)
	var categories domain.Category
	err := row.Scan(
		&categories.ID,
		&categories.Name,
		&categories.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Category{}, domain.NewNotFoundError(fmt.Sprintf("category with name %s not found", name), err)
		}
		return domain.Category{}, domain.NewInternalError(err.Error(), err)
	}

	return categories, nil
}

func (c *categoryRepository) Update(ctx context.Context, tx *sql.Tx, category domain.Category) (int64, error) {
	res, err := tx.ExecContext(
		ctx,
		updateCategoryQuery,
		category.Name,
		category.ID,
	)
	if err != nil {
		return 0, domain.NewInternalError(err.Error(), err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, domain.NewInternalError("fail to get rows affected", err)
	}

	if rowsAffected == 0 {
		return 0, domain.NewNotFoundError(fmt.Sprintf("category with id %d not found", category.ID), err)
	}

	return rowsAffected, nil
}

func (c *categoryRepository) Delete(ctx context.Context, tx *sql.Tx, id int64) (int64, error) {
	res, err := tx.ExecContext(
		ctx,
		deleteCategoryQuery,
		id,
	)
	if err != nil {
		return 0, domain.NewInternalError(err.Error(), err)
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, domain.NewInternalError("fail to get last insert id", err)
	}
	return id, nil
}
