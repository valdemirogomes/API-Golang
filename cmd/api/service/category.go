package service

import (
	"context"
	"database/sql"

	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/config"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/domain"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/dto"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/lib/helperdb"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/lib/mysql"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/repository"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, category domain.Category) (dto.CategoryResponseDTO, error)
	GetCategories(ctx context.Context, params dto.SearchParams) (dto.CategoryListResponseDTO, error)
	UpdateCategory(ctx context.Context, category dto.CategoryDTO, id int64) (dto.CategoryResponseDTO, error)
	DeleteCategory(ctx context.Context, id int64) (int64, error)
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
	db                 mysql.DB
	config             config.Environment
}

func NewCategoryService(categoryRepository repository.CategoryRepository, db mysql.DB, config config.Environment) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
		db:                 db,
		config:             config,
	}
}

func (c *categoryService) CreateCategory(ctx context.Context, category domain.Category) (dto.CategoryResponseDTO, error) {
	var categoryDomain domain.Category
	var id int64
	var err error

	if category.Name == "" {
		return dto.CategoryResponseDTO{}, domain.NewBadRequest("category name is required", nil)
	}

	txErr := c.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		id, err = c.categoryRepository.Create(ctx, tx, category)
		if err != nil {
			return err
		}
		categoryDomain, err = c.categoryRepository.FindByID(ctx, tx, id)
		return err
	})
	if txErr != nil {
		return dto.CategoryResponseDTO{}, txErr
	}
	return dto.CategoryResponseDTO{
		ID:   categoryDomain.ID,
		Name: categoryDomain.Name,
	}, nil
}

func (c *categoryService) GetCategories(ctx context.Context, param dto.SearchParams) (dto.CategoryListResponseDTO, error) {
	var categoriesResponse dto.CategoryListResponseDTO
	txErr := c.db.WithoutTransaction(ctx, func(tx helperdb.Tx) error {
		categories, total, err := c.categoryRepository.FindAll(ctx, tx, param)
		if err != nil {
			return err
		}
		var categoriesDTO []dto.CategoryResponseDTO
		for _, category := range categories {
			categoriesDTO = append(categoriesDTO, dto.CategoryResponseDTO{
				ID:   category.ID,
				Name: category.Name,
			})
		}

		categoriesResponse = dto.CategoryListResponseDTO{
			Data: categoriesDTO,
			Metadata: dto.Metadata{
				Total:        total,
				Limit:        *param.Limit,
				Offset:       *param.Offset,
				TotalEntries: int64(len(categoriesDTO)),
			},
		}
		return nil
	})
	if txErr != nil {
		return dto.CategoryListResponseDTO{}, txErr
	}
	return categoriesResponse, nil
}

func (c *categoryService) UpdateCategory(ctx context.Context, category dto.CategoryDTO, id int64) (dto.CategoryResponseDTO, error) {
	var categoryDTO dto.CategoryResponseDTO
	txErr := c.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		categoryDomain, err := c.categoryRepository.FindByID(ctx, tx, id)
		if err != nil {
			return err
		}
		categoryUpload := domain.Category{
			ID:   categoryDomain.ID,
			Name: category.Name,
		}

		_, err = c.categoryRepository.Update(ctx, tx, categoryUpload)
		if err != nil {
			return err
		}

		categoryDTO = dto.CategoryResponseDTO{
			ID:   categoryUpload.ID,
			Name: categoryUpload.Name,
		}

		return nil
	})
	if txErr != nil {
		return dto.CategoryResponseDTO{}, txErr
	}
	return categoryDTO, nil
}

func (c *categoryService) DeleteCategory(ctx context.Context, id int64) (int64, error) {
	var err error
	txErr := c.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err = c.categoryRepository.Delete(ctx, tx, id)
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return 0, txErr
	}
	return id, nil
}
