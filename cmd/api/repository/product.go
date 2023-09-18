package repository

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/domain"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/dto"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/lib/helperdb"
)

type ProductRepository interface {
	Create(ctx context.Context, tx *sql.Tx, product domain.Product) (int64, error)
	FindAll(ctx context.Context, tx helperdb.Tx, params dto.SearchParams) ([]domain.Product, int64, error)
	FindByID(ctx context.Context, tx helperdb.Tx, id int64) (domain.Product, error)
	FindByCategory(ctx context.Context, tx helperdb.Tx, categoryID int64) ([]domain.Product, error)
	Update(ctx context.Context, tx *sql.Tx, product domain.Product) (int64, error)
	Delete(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
}

type productRepository struct {
}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

const (
	createProductQuery   = "INSERT INTO products (title, description, price, image, created_at, category_id) VALUES ( ?, ?, ?, ?, NOW(), ?)"
	findByIDProductQuery = "SELECT id, title, description, price, image, created_at, category_id FROM products WHERE id = ?"
	findAllProductsQuery = "SELECT id, title, description, price, image, created_at, category_id FROM products"
	findByCategoryQuery  = "SELECT id, title, description, price, image, created_at, category_id FROM products WHERE category_id = ?"
	updateProductQuery   = "UPDATE products SET title = ?, description = ?, price = ?, image = ?, category_id = ? WHERE id = ?"
	deleteProductQuery   = "DELETE FROM products WHERE id = ?"
)

func (p *productRepository) Create(ctx context.Context, tx *sql.Tx, product domain.Product) (int64, error) {
	res, err := tx.ExecContext(
		ctx,
		createProductQuery,
		product.Title,
		product.Description,
		product.Price,
		product.Image,
		product.CategoryID,
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

func (p *productRepository) FindAll(ctx context.Context, tx helperdb.Tx, params dto.SearchParams) ([]domain.Product, int64, error) {
	query, queryParams := GetSearchQuery(params, findAllProductsQuery)
	var products []domain.Product
	var total int64
	err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM products").Scan(&total)
	if err != nil {
		return nil, 0, domain.NewInternalError(err.Error(), err)
	}
	rows, err := tx.QueryContext(ctx, query.String(), queryParams...)
	if err != nil {
		return nil, 0, domain.NewInternalError("fail to execute query", err)
	}
	defer rows.Close()
	for rows.Next() {
		var product domain.Product
		err = rows.Scan(
			&product.ID,
			&product.Title,
			&product.Description,
			&product.Price,
			&product.Image,
			&product.CreatedAt,
			&product.CategoryID,
		)
		if err != nil {
			return nil, 0, domain.NewInternalError("fail to scan row", err)
		}
		products = append(products, product)
	}
	if len(products) == 0 {
		return nil, 0, domain.NewNotFoundError(fmt.Sprintf("product not found"), errors.New("product not found"))
	}
	return products, total, nil
}

func (p *productRepository) FindByID(ctx context.Context, tx helperdb.Tx, id int64) (domain.Product, error) {
	row := tx.QueryRowContext(ctx, findByIDProductQuery, id)
	var product domain.Product
	err := row.Scan(
		&product.ID,
		&product.Title,
		&product.Description,
		&product.Price,
		&product.Image,
		&product.CreatedAt,
		&product.CategoryID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, domain.NewNotFoundError(fmt.Sprintf("product with ID %d not found", id), err)
		}
		return product, domain.NewInternalError(err.Error(), err)
	}
	return product, nil
}

func (p *productRepository) FindByCategory(ctx context.Context, tx helperdb.Tx, categoryID int64) ([]domain.Product, error) {
	rows, err := tx.QueryContext(ctx, findByCategoryQuery, categoryID)
	if err != nil {
		return nil, domain.NewInternalError(err.Error(), err)
	}
	defer rows.Close()
	var products []domain.Product

	isNotFound := true

	for rows.Next() {
		var product domain.Product
		err = rows.Scan(
			&product.ID,
			&product.Title,
			&product.Description,
			&product.Price,
			&product.Image,
			&product.CreatedAt,
			&product.CategoryID,
		)
		if err != nil {
			return nil, domain.NewInternalError("fail to scan row", err)
		}
		products = append(products, product)
		isNotFound = false
	}

	if isNotFound {
		return nil, domain.NewNotFoundError(fmt.Sprintf("product with category ID %d not found", categoryID), err)
	}

	return products, nil
}

func (p *productRepository) Update(ctx context.Context, tx *sql.Tx, product domain.Product) (int64, error) {
	res, err := tx.ExecContext(
		ctx,
		updateProductQuery,
		product.Title,
		product.Description,
		product.Price,
		product.Image,
		product.CategoryID,
		product.ID,
	)
	if err != nil {
		return 0, domain.NewInternalError(err.Error(), err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, domain.NewInternalError("fail to get rows affected", err)
	}
	if rowsAffected == 0 {
		return 0, domain.NewNotFoundError(fmt.Sprintf("product with ID %d not found", product.ID), err)
	}
	return rowsAffected, nil
}

func (p *productRepository) Delete(ctx context.Context, tx *sql.Tx, id int64) (int64, error) {
	res, err := tx.ExecContext(ctx, deleteProductQuery, id)
	if err != nil {
		return 0, domain.NewInternalError(err.Error(), err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, domain.NewInternalError("fail to get rows affected", err)
	}
	if rowsAffected == 0 {
		return 0, domain.NewNotFoundError(fmt.Sprintf("product with ID %d not found", id), err)
	}
	return rowsAffected, nil
}

func GetSearchQuery(params dto.SearchParams, querySQL string) (bytes.Buffer, []interface{}) {
	var query bytes.Buffer
	query.WriteString(querySQL)
	query.WriteString(" WHERE 1=1")
	queryParams := make([]interface{}, 0)

	if params.Min != nil && params.Max != nil {
		query.WriteString(" AND price BETWEEN ? AND ?")
		queryParams = append(queryParams, *params.Min, *params.Max)
	}

	if params.Sort == nil {
		query.WriteString(" ORDER BY created_at DESC")
	}

	if querySQL == findAllProductsQuery {
		if params.Name != nil {
			query.WriteString(" AND title LIKE ?")
			queryParams = append(queryParams, "%"+*params.Name+"%")
		}
		if params.Title != nil {
			query.WriteString(" AND title LIKE ?")
			queryParams = append(queryParams, "%"+*params.Title+"%")
		}
		if params.Sort != nil {
			query.WriteString(" ORDER BY title ")
			query.WriteString(*params.Sort)
		}

	}

	if querySQL == findAllCategoryQuery {
		if params.Name != nil {
			query.WriteString(" AND name LIKE ?")
			queryParams = append(queryParams, "%"+*params.Name+"%")
		}
		if params.Title != nil {
			query.WriteString(" AND name LIKE ?")
			queryParams = append(queryParams, "%"+*params.Title+"%")
		}

		if params.Sort != nil {
			query.WriteString(" ORDER BY name ")
			query.WriteString(*params.Sort)
		}
	}

	if params.Limit != nil {
		query.WriteString(" LIMIT ?")
		queryParams = append(queryParams, *params.Limit)
	}

	if params.Offset != nil {
		query.WriteString(" OFFSET ?")
		offset := (*params.Offset - 1) * *params.Limit
		queryParams = append(queryParams, offset)
	}

	return query, queryParams
}
