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

type ProductService interface {
	GetProducts(ctx context.Context, params dto.SearchParams) (dto.ProductResponse, error)
	FindById(ctx context.Context, id int64) (dto.ProductDTO, error)
	CreateProduct(ctx context.Context, productDTO dto.ProductDTO) (dto.ProductDTO, error)
	GetProductsByCategory(ctx context.Context, category string) ([]dto.ProductDTO, error)
	UpdateProduct(ctx context.Context, product dto.ProductUpdateDTO, id int64) (dto.ProductDTO, error)
	DeleteProduct(ctx context.Context, id int64) error
}

type productService struct {
	productRepository  repository.ProductRepository
	categoryRepository repository.CategoryRepository
	db                 mysql.DB
	env                config.Environment
}

func NewProductService(productRepository repository.ProductRepository, categoryRepository repository.CategoryRepository, db mysql.DB, env config.Environment) *productService {
	return &productService{
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
		db:                 db,
		env:                env,
	}
}

func (p *productService) GetProducts(ctx context.Context, param dto.SearchParams) (dto.ProductResponse, error) {
	var productsResponse dto.ProductResponse

	txErr := p.db.WithoutTransaction(ctx, func(tx helperdb.Tx) error {
		products, total, err := p.productRepository.FindAll(ctx, tx, param)
		if err != nil {
			return err
		}

		var productsDTO []dto.ProductDTO
		for _, product := range products {
			categoryName, err := p.categoryRepository.FindByID(ctx, tx, product.CategoryID)
			if err != nil {
				return err
			}
			productsDTO = append(productsDTO, dto.ProductDTO{
				ID:          product.ID,
				Title:       product.Title,
				Description: product.Description,
				Price:       product.Price,
				Image:       product.Image,
				Category:    categoryName.Name,
			})
		}

		productsResponse = dto.ProductResponse{
			Data: productsDTO,
			Metadata: dto.Metadata{
				Total:        total,
				Limit:        *param.Limit,
				Offset:       *param.Offset,
				TotalEntries: int64(len(productsDTO)),
			},
		}

		return nil
	})
	if txErr != nil {
		return dto.ProductResponse{}, txErr
	}
	return productsResponse, nil
}

func (p *productService) FindById(ctx context.Context, id int64) (dto.ProductDTO, error) {
	var categoryDomain domain.Category
	var productDomain domain.Product
	var err error

	txErr := p.db.WithoutTransaction(ctx, func(tx helperdb.Tx) error {
		productDomain, err = p.productRepository.FindByID(ctx, tx, id)
		if err != nil {
			return err
		}
		categoryDomain, err = p.categoryRepository.FindByID(ctx, tx, productDomain.CategoryID)
		return err
	})
	if txErr != nil {
		return dto.ProductDTO{}, txErr
	}
	return dto.ProductDTO{
		ID:          productDomain.ID,
		Title:       productDomain.Title,
		Description: productDomain.Description,
		Price:       productDomain.Price,
		Image:       productDomain.Image,
		Category:    categoryDomain.Name,
	}, nil

}

func (p *productService) CreateProduct(ctx context.Context, productDTO dto.ProductDTO) (dto.ProductDTO, error) {
	var id int64
	txErr := p.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		category, err := p.categoryRepository.FindByName(ctx, tx, productDTO.Category)
		if err != nil {
			return err
		}
		productDomain := domain.Product{
			Title:       productDTO.Title,
			Description: productDTO.Description,
			Price:       productDTO.Price,
			Image:       productDTO.Image,
			CategoryID:  category.ID,
		}
		id, err = p.productRepository.Create(ctx, tx, productDomain)
		return err
	})
	if txErr != nil {
		return dto.ProductDTO{}, txErr
	}
	productDTO.ID = id
	return productDTO, nil
}

func (p *productService) GetProductsByCategory(ctx context.Context, category string) ([]dto.ProductDTO, error) {
	var productsDTO []dto.ProductDTO
	txErr := p.db.WithoutTransaction(ctx, func(tx helperdb.Tx) error {
		categoryDomain, err := p.categoryRepository.FindByName(ctx, tx, category)
		if err != nil {
			return err
		}
		products, err := p.productRepository.FindByCategory(ctx, tx, categoryDomain.ID)
		if err != nil {
			return err
		}
		for _, productDomain := range products {
			productsDTO = append(productsDTO, dto.ProductDTO{
				ID:          productDomain.ID,
				Title:       productDomain.Title,
				Description: productDomain.Description,
				Price:       productDomain.Price,
				Image:       productDomain.Image,
				Category:    categoryDomain.Name,
			})
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}
	return productsDTO, nil
}

func (p *productService) UpdateProduct(ctx context.Context, product dto.ProductUpdateDTO, id int64) (dto.ProductDTO, error) {
	var productDTO dto.ProductDTO
	txErr := p.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		productDomain, err := p.productRepository.FindByID(ctx, tx, id)
		if err != nil {
			return err
		}
		categoryName, err := p.categoryRepository.FindByID(ctx, tx, productDomain.CategoryID)
		if err != nil {
			return err
		}

		productUpdate := domain.Product{
			ID:          productDomain.ID,
			Title:       product.Title,
			Description: product.Description,
			Price:       product.Price,
			Image:       product.Image,
			CategoryID:  productDomain.CategoryID,
		}

		_, err = p.productRepository.Update(ctx, tx, productUpdate)
		if err != nil {
			return err
		}

		productDTO = dto.ProductDTO{
			Title:       productUpdate.Title,
			Description: productUpdate.Description,
			Price:       productUpdate.Price,
			Image:       productUpdate.Image,
			Category:    categoryName.Name,
		}

		return nil
	})
	if txErr != nil {
		return dto.ProductDTO{}, txErr
	}
	return productDTO, nil
}

func (p *productService) DeleteProduct(ctx context.Context, id int64) error {
	txErr := p.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := p.productRepository.Delete(ctx, tx, id)
		return err
	})
	if txErr != nil {
		return txErr
	}
	return nil
}
