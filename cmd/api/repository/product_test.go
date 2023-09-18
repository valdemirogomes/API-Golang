package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/domain"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/dto"
	"github.com/stretchr/testify/assert"
)

var (
	// once   sync.Once
	// test   *dbtest.Test
	limit  = int64(10)
	offset = int64(1)
	min    = float64(0)
	max    = float64(100)
	sort   = "ASC"
	title  = "test"
	name   = "test"
)
var productRows = []string{
	"id",
	"title",
	"description",
	"price",
	"image",
	"created_at",
	"category_id",
}

func InitialCommonMocks() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()

	return db, mock

}

func InitialMockDBProduct() domain.Product {
	categoryID := int64(1)
	return domain.Product{
		ID:          1,
		Title:       "Test Product",
		Description: "Test Description",
		Price:       10.00,
		Image:       "test.jpg",
		CategoryID:  categoryID,
		CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}
func TestCreateProduct_IntoTx_WithoutError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	mock.ExpectBegin()
	query := QueryReplace(createProductQuery)

	mock.ExpectExec(query).WithArgs(
		product.Title,
		product.Description,
		product.Price,
		product.Image,
		product.CategoryID,
	).WillReturnResult(sqlmock.NewResult(product.ID, 1))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	id, err := repo.Create(context.Background(), tx, product)

	assert.NoError(t, err, "Error should not be returned")

	assert.Equal(t, product.ID, id)

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestCreateProduct_IntoTx_WithError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	mock.ExpectBegin()
	query := QueryReplace(createProductQuery)

	mock.ExpectExec(query).WithArgs(
		product.Title,
		product.Description,
		product.Price,
		product.Image,
		product.CategoryID,
	).WillReturnError(sql.ErrConnDone)

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	id, err := repo.Create(context.Background(), tx, product)

	assert.Error(t, err, "Error should be returned")

	assert.Equal(t, int64(0), id)

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestCreateProduct_IntoTx_WithErrorInLastInsertId(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	mock.ExpectBegin()
	query := QueryReplace(createProductQuery)

	mock.ExpectExec(query).WithArgs(
		product.Title,
		product.Description,
		product.Price,
		product.Image,
		product.CategoryID,
	).WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Create(context.Background(), tx, product)
	assert.NotNil(t, err, "Error should be returned")
}

func TestFindAllProducts(t *testing.T) {
	db, mock := InitialCommonMocks()

	defer db.Close()

	products := []domain.Product{
		InitialMockDBProduct(),
		InitialMockDBProduct(),
		InitialMockDBProduct(),
	}

	query := QueryReplace(findAllProductsQuery)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery(query).
		WillReturnRows(mock.NewRows(productRows).
			AddRow(products[0].ID, products[0].Title, products[0].Description, products[0].Price, products[0].Image, products[0].CreatedAt, products[0].CategoryID).
			AddRow(products[1].ID, products[1].Title, products[1].Description, products[1].Price, products[1].Image, products[1].CreatedAt, products[1].CategoryID).
			AddRow(products[2].ID, products[2].Title, products[2].Description, products[2].Price, products[2].Image, products[2].CreatedAt, products[2].CategoryID))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	params := dto.SearchParams{
		Limit:  &limit,
		Offset: &offset,
		Min:    &min,
		Max:    &max,
		Sort:   &sort,
		Title:  &title,
	}

	result, total, err := repo.FindAll(context.Background(), tx, params)
	assert.NoError(t, err, "Error should not be returned")
	assert.Equal(t, len(products), len(result), "Number of returned products should match")
	assert.Equal(t, int64(len(products)), total, "Total count should match")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestFindAllProducts_WithoutError_WithName(t *testing.T) {
	db, mock := InitialCommonMocks()

	defer db.Close()

	products := []domain.Product{
		InitialMockDBProduct(),
		InitialMockDBProduct(),
		InitialMockDBProduct(),
	}

	query := QueryReplace(findAllProductsQuery)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery(query).
		WillReturnRows(mock.NewRows(productRows).
			AddRow(products[0].ID, products[0].Title, products[0].Description, products[0].Price, products[0].Image, products[0].CreatedAt, products[0].CategoryID).
			AddRow(products[1].ID, products[1].Title, products[1].Description, products[1].Price, products[1].Image, products[1].CreatedAt, products[1].CategoryID).
			AddRow(products[2].ID, products[2].Title, products[2].Description, products[2].Price, products[2].Image, products[2].CreatedAt, products[2].CategoryID))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	params := dto.SearchParams{
		Limit:  &limit,
		Offset: &offset,
		Min:    &min,
		Max:    &max,
		Sort:   &sort,
		Name:   &name,
	}

	result, total, err := repo.FindAll(context.Background(), tx, params)
	assert.NoError(t, err, "Error should not be returned")
	assert.Equal(t, len(products), len(result), "Number of returned products should match")
	assert.Equal(t, int64(len(products)), total, "Total count should match")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestFindAllProducts_WithErrorInScan(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	query := QueryReplace(findAllProductsQuery)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery(query).
		WillReturnRows(mock.NewRows([]string{"id"}).
			AddRow(1))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	params := dto.SearchParams{
		Limit:  &limit,
		Offset: &offset,
		Min:    nil,
		Max:    nil,
		Sort:   &sort,
		Title:  nil,
	}

	_, _, err = repo.FindAll(context.Background(), tx, params)

	assert.Error(t, err, "Error should be returned during scanning")
	assert.Contains(t, err.Error(), "fail to scan row", "Error message should contain expected text")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")

}

func TestFindAllProducts_WithErrorInCount(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	query := QueryReplace(findAllProductsQuery)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
		WillReturnError(errors.New("fail to execute query"))
	mock.ExpectQuery(query).
		WillReturnRows(mock.NewRows(productRows))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	params := dto.SearchParams{
		Limit:  nil,
		Offset: nil,
		Min:    nil,
		Max:    nil,
		Sort:   nil,
		Title:  nil,
	}

	_, _, err = repo.FindAll(context.Background(), tx, params)

	assert.Error(t, err, "Error should be returned during scanning")
	assert.Contains(t, err.Error(), "fail to execute query", "Error message should contain expected text")

}

func TestFindAllProducts_WithErrorNotFound(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	query := QueryReplace(findAllProductsQuery)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(query).
		WillReturnRows(mock.NewRows(productRows))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	params := dto.SearchParams{
		Limit:  nil,
		Offset: nil,
		Min:    nil,
		Max:    nil,
		Sort:   nil,
		Title:  nil,
	}

	_, _, err = repo.FindAll(context.Background(), tx, params)

	assert.Error(t, err, "Error should be returned during scanning")
	assert.Contains(t, err.Error(), "product not found", "Error message should contain expected text")

}

func TestFindAll_WithError_InQuery(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	query := QueryReplace(findAllProductsQuery)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery(query).
		WillReturnError(errors.New("some error"))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	params := dto.SearchParams{
		Limit:  nil,
		Offset: nil,
		Min:    nil,
		Max:    nil,
		Sort:   nil,
		Title:  nil,
	}

	_, _, err = repo.FindAll(context.Background(), tx, params)

	assert.Error(t, err, "Error should be returned during scanning")
	assert.Contains(t, err.Error(), "fail to execute query", "Error message should contain expected text")

}

func TestFindByID_WithoutError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	query := QueryReplace(findByIDProductQuery)

	mock.ExpectBegin()
	mock.ExpectQuery(query).
		WithArgs(product.ID).
		WillReturnRows(mock.NewRows(productRows).
			AddRow(product.ID, product.Title, product.Description, product.Price, product.Image, product.CreatedAt, product.CategoryID))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	result, err := repo.FindByID(context.Background(), tx, product.ID)
	assert.NoError(t, err, "Error should not be returned")
	assert.Equal(t, product, result, "Product should match")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestFindByID_WithError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	query := QueryReplace(findByIDProductQuery)

	mock.ExpectBegin()
	mock.ExpectQuery(query).
		WithArgs(product.ID).
		WillReturnError(sql.ErrConnDone)

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.FindByID(context.Background(), tx, product.ID)
	assert.Error(t, err, "Error should be returned")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestFindByID_WithErrorErrNoRows(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	query := QueryReplace(findByIDProductQuery)

	mock.ExpectBegin()
	mock.ExpectQuery(query).
		WithArgs(product.ID).
		WillReturnError(sql.ErrNoRows)

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.FindByID(context.Background(), tx, product.ID)
	assert.Error(t, err, "Error should be returned")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestFindByCategory_WithResults(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	categoryID := int64(123)

	query := QueryReplace(findByCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery(query).
		WithArgs(categoryID).
		WillReturnRows(mock.NewRows(productRows).
			AddRow(1, "Product 1", "Description 1", 10.00, "image1.jpg", time.Now(), categoryID).
			AddRow(2, "Product 2", "Description 2", 20.00, "image2.jpg", time.Now(), categoryID))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	result, err := repo.FindByCategory(context.Background(), tx, categoryID)
	assert.NoError(t, err, "Error should not be returned")
	assert.Len(t, result, 2, "Number of returned products should be 2")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestFindByCategory_QueryError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	categoryID := int64(123)

	query := QueryReplace(findByCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery(query).
		WithArgs(categoryID).
		WillReturnError(errors.New("some error"))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.FindByCategory(context.Background(), tx, categoryID)
	assert.Error(t, err, "Error should be returned")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestFindByCategory_NoRowsError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	categoryID := int64(123)

	query := QueryReplace(findByCategoryQuery)

	mock.ExpectQuery(query).WithArgs(categoryID).WillReturnError(sql.ErrNoRows)

	repo := NewProductRepository()

	_, err := repo.FindByCategory(context.Background(), db, categoryID)
	assert.Error(t, err, "Error should be returned")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestFindByCategory_ScanError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	categoryID := int64(123)

	query := QueryReplace(findByCategoryQuery)

	mock.ExpectQuery(query).
		WithArgs(categoryID).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1))

	repo := NewProductRepository()

	_, err := repo.FindByCategory(context.Background(), db, categoryID)
	assert.Error(t, err, "Error should be returned")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestFindByCategory_isNotFound(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	categoryID := int64(123)

	query := QueryReplace(findByCategoryQuery)

	mock.ExpectQuery(query).
		WithArgs(categoryID).
		WillReturnRows(mock.NewRows(productRows))

	repo := NewProductRepository()

	_, err := repo.FindByCategory(context.Background(), db, categoryID)
	assert.Error(t, err, "Error should be returned")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestUpdateProduct_WithoutError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	query := QueryReplace(updateProductQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(product.Title, product.Description, product.Price, product.Image, product.CategoryID, product.ID).
		WillReturnResult(sqlmock.NewResult(product.ID, 1))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Update(context.Background(), tx, product)
	assert.NoError(t, err, "Error should not be returned")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestUpdateProduct_WithError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	query := QueryReplace(updateProductQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(product.Title, product.Description, product.Price, product.Image, product.CategoryID, product.ID).
		WillReturnError(sql.ErrConnDone)

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Update(context.Background(), tx, product)
	assert.Error(t, err, "Error should be returned")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestUpdateProduct_WithErrorInResult(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	query := QueryReplace(updateProductQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(product.Title, product.Description, product.Price, product.Image, product.CategoryID, product.ID).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Update(context.Background(), tx, product)
	assert.NotNil(t, err, "Error should be returned")
}

func TestUpdateProduct_WithZeroRowsAffected(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	query := QueryReplace(updateProductQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(product.Title, product.Description, product.Price, product.Image, product.CategoryID, product.ID).
		WillReturnResult(sqlmock.NewResult(product.ID, 0))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Update(context.Background(), tx, product)
	assert.NotNil(t, err, "Error should be returned")
}

func TestDeleteProduct_WithoutError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	query := QueryReplace(deleteProductQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(product.ID).
		WillReturnResult(sqlmock.NewResult(product.ID, 1))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Delete(context.Background(), tx, product.ID)
	assert.NoError(t, err, "Error should not be returned")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestDeleteProduct_WithError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	query := QueryReplace(deleteProductQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(product.ID).
		WillReturnError(sql.ErrConnDone)

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Delete(context.Background(), tx, product.ID)
	assert.Error(t, err, "Error should be returned")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestDeleteProduct_WithErrorInResult(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	query := QueryReplace(deleteProductQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(product.ID).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Delete(context.Background(), tx, product.ID)
	assert.NotNil(t, err, "Error should be returned")
}

func TestDeleteProduct_WithZeroRowsAffected(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	product := InitialMockDBProduct()

	query := QueryReplace(deleteProductQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(product.ID).
		WillReturnResult(sqlmock.NewResult(product.ID, 0))

	repo := NewProductRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Delete(context.Background(), tx, product.ID)
	assert.NotNil(t, err, "Error should be returned")
}

func QueryReplace(query string) string {
	replacements := map[string]string{
		"?": "\\?",
		")": "\\)",
		"(": "\\(",
	}
	for k, v := range replacements {
		query = strings.ReplaceAll(query, k, v)
	}
	return query
}
