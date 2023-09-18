package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/dto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/domain"
	"github.com/stretchr/testify/assert"
)

var categoryRows = []string{
	"id",
	"name",
	"created_at",
}

func InitialMockDBCategory() domain.Category {
	return domain.Category{
		ID:        1,
		Name:      "test",
		CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}
func TestCreateCategory_IntoTx_WithoutError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	mock.ExpectBegin()
	query := QueryReplace(createCategoryQuery)

	mock.ExpectExec(query).WithArgs(
		category.Name,
	).WillReturnResult(sqlmock.NewResult(category.ID, 1))

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	id, err := repo.Create(context.Background(), tx, category)

	assert.NoError(t, err, "Error should not be returned")

	assert.Equal(t, category.ID, id)

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestCreateCategory_IntoTx_WithError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	mock.ExpectBegin()
	query := QueryReplace(createCategoryQuery)

	mock.ExpectExec(query).WithArgs(
		category.Name,
	).WillReturnError(sql.ErrConnDone)

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Create(context.Background(), tx, category)
	assert.Error(t, err, "Error should be returned")

	assert.NoError(t, mock.ExpectationsWereMet(), "Error should not be returned")
}

func TestCreateCategory_IntoTx_WithErrorInLastInsertId(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	mock.ExpectBegin()
	query := QueryReplace(createCategoryQuery)

	mock.ExpectExec(query).WithArgs(
		category.Name,
	).WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Create(context.Background(), tx, category)
	assert.NotNil(t, err, "Error should be returned")
}

func TestFindAll_WithoutError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := []domain.Category{
		InitialMockDBCategory(),
		InitialMockDBCategory(),
	}

	query := QueryReplace(findAllCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM categories").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(category)))
	mock.ExpectQuery(query).
		WillReturnRows(mock.NewRows(categoryRows).
			AddRow(category[0].ID, category[0].Name, category[0].CreatedAt).
			AddRow(category[1].ID, category[1].Name, category[1].CreatedAt))

	repo := NewCategoryRepository()

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

	categories, total, err := repo.FindAll(context.Background(), tx, params)

	assert.NoError(t, err, "Error should not be returned")
	assert.Equal(t, len(category), len(categories), "Number of returned products should match")
	assert.Equal(t, int64(len(category)), total, "Total count should match")

}

func TestFindAll_WithTitle(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := []domain.Category{
		InitialMockDBCategory(),
		InitialMockDBCategory(),
	}

	query := QueryReplace(findAllCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM categories").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(category)))
	mock.ExpectQuery(query).
		WillReturnRows(mock.NewRows(categoryRows).
			AddRow(category[0].ID, category[0].Name, category[0].CreatedAt).
			AddRow(category[1].ID, category[1].Name, category[1].CreatedAt))

	repo := NewCategoryRepository()

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

	categories, total, err := repo.FindAll(context.Background(), tx, params)

	assert.NoError(t, err, "Error should not be returned")
	assert.Equal(t, len(category), len(categories), "Number of returned products should match")
	assert.Equal(t, int64(len(category)), total, "Total count should match")
}

func TestFindAll_WithError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	query := QueryReplace(findAllCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM categories").
		WillReturnError(sql.ErrConnDone)
	mock.ExpectQuery(query).
		WillReturnError(sql.ErrConnDone)

	repo := NewCategoryRepository()

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

	_, _, err = repo.FindAll(context.Background(), tx, params)

	assert.Error(t, err, "Error should be returned")
}

func TestFindAll_WithErrorInScan(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	query := QueryReplace(findAllCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM categories").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery(query).
		WillReturnRows(mock.NewRows([]string{"id"}).
			AddRow(1))

	repo := NewCategoryRepository()

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

	_, _, err = repo.FindAll(context.Background(), tx, params)

	assert.Error(t, err, "Error should be returned")
}

func TestFindAll_WithErrorInCount(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	query := QueryReplace(findAllCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM categories").
		WillReturnError(errors.New("some error"))
	mock.ExpectQuery(query).
		WillReturnRows(mock.NewRows(categoryRows))

	repo := NewCategoryRepository()

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

}

func TestFindAll_WithErrorInQuery(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	query := QueryReplace(findAllCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM categories").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery(query).
		WillReturnError(errors.New("some error"))

	repo := NewCategoryRepository()

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

}

func TestFindCategoryByID_WithoutError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(findByIDCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery(query).
		WithArgs(category.ID).
		WillReturnRows(mock.NewRows(categoryRows).
			AddRow(category.ID, category.Name, category.CreatedAt))

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	categoryResult, err := repo.FindByID(context.Background(), tx, category.ID)

	assert.NoError(t, err, "Error should not be returned")
	assert.Equal(t, category, categoryResult, "Category should match")

}

func TestFindCategoryByID_WithError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(findByIDCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery(query).
		WithArgs(category.ID).
		WillReturnError(sql.ErrConnDone)

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.FindByID(context.Background(), tx, category.ID)

	assert.Error(t, err, "Error should be returned")
}

func TestFindCategoryByID_WithErrorErrNoRows(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(findByIDCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery(query).
		WithArgs(category.ID).
		WillReturnError(sql.ErrNoRows)

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.FindByID(context.Background(), tx, category.ID)

	assert.Error(t, err, "Error should be returned")
}

func TestFindByName_WithoutError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(findByNameCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery(query).
		WithArgs(category.Name).
		WillReturnRows(mock.NewRows(categoryRows).
			AddRow(category.ID, category.Name, category.CreatedAt))

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	categoryResult, err := repo.FindByName(context.Background(), tx, category.Name)

	assert.NoError(t, err, "Error should not be returned")
	assert.Equal(t, category, categoryResult, "Category should match")

}

func TestFindByName_WithError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(findByNameCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery(query).
		WithArgs(category.Name).
		WillReturnError(sql.ErrConnDone)

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.FindByName(context.Background(), tx, category.Name)

	assert.Error(t, err, "Error should be returned")
}

func TestFindByName_WithErrorErrNoRows(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(findByNameCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectQuery(query).
		WithArgs(category.Name).
		WillReturnError(sql.ErrNoRows)

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.FindByName(context.Background(), tx, category.Name)

	assert.Error(t, err, "Error should be returned")
}

func TestUpdate_WithoutError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(updateCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(category.Name, category.ID).
		WillReturnResult(sqlmock.NewResult(category.ID, 1))

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	id, err := repo.Update(context.Background(), tx, category)

	assert.NoError(t, err, "Error should not be returned")
	assert.Equal(t, category.ID, id, "Category should match")

}

func TestUpdate_WithError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(updateCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(category.Name, category.ID).
		WillReturnError(sql.ErrConnDone)

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Update(context.Background(), tx, category)

	assert.Error(t, err, "Error should be returned")
}

func TestUpdate_WithErrorErrNoRows(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(updateCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(category.Name, category.ID).
		WillReturnResult(sqlmock.NewResult(category.ID, 0))

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Update(context.Background(), tx, category)

	assert.Error(t, err, "Error should be returned")
}

func TestUpdate_WithErroRowsAffected(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(updateCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(category.Name, category.ID).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Update(context.Background(), tx, category)

	assert.Error(t, err, "Error should be returned")
}

func TestDelete_WithoutError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(deleteCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(category.ID).
		WillReturnResult(sqlmock.NewResult(category.ID, 1))

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	id, err := repo.Delete(context.Background(), tx, category.ID)

	assert.NoError(t, err, "Error should not be returned")
	assert.Equal(t, category.ID, id, "Category should match")

}

func TestDelete_WithError(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(deleteCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(category.ID).
		WillReturnError(sql.ErrConnDone)

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Delete(context.Background(), tx, category.ID)

	assert.Error(t, err, "Error should be returned")
}

func TestDelete_WithErrorLastInsertId(t *testing.T) {
	db, mock := InitialCommonMocks()
	defer db.Close()

	category := InitialMockDBCategory()

	query := QueryReplace(deleteCategoryQuery)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(category.ID).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))

	repo := NewCategoryRepository()

	tx, err := db.Begin()
	assert.NoError(t, err, "Error should not be returned")

	_, err = repo.Delete(context.Background(), tx, category.ID)

	assert.Error(t, err, "Error should be returned")
}
