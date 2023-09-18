package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/domain"
	"github.com/mercadolibre/fury_go-core/pkg/log"

	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/config"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/dto"
	errorhandling "github.com/melisource/fury_go-dev-base-3-v2/cmd/api/lib/error_handling"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/service"
	"github.com/mercadolibre/fury_go-core/pkg/web"
)

type ProductController interface {
	HandleGetProducts(w http.ResponseWriter, r *http.Request) error
	HandleGetProductByID(w http.ResponseWriter, r *http.Request) error
	HandleCreateProduct(w http.ResponseWriter, r *http.Request) error
	HandleFindProductByCategory(w http.ResponseWriter, r *http.Request) error
	HandleUpdateProduct(w http.ResponseWriter, r *http.Request) error
	HandleDeleteProduct(w http.ResponseWriter, r *http.Request) error
}

type productController struct {
	productService service.ProductService
	config         config.Environment
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func NewProductController(productService service.ProductService, config config.Environment) ProductController {
	return &productController{
		productService: productService,
		config:         config,
	}
}

// HandleGetProducts godoc
// @Summary Get products
// @Description Get products
// @Tags products
// @Accept  json
// @Produce  json
// @Param sort query string false "sort"
// @Param title query string false "title"
// @Param name query string false "name"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Param min query float64 false "min"
// @Param max query float64 false "max"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} ErrorMessage
// @Failure 400 {object} ErrorMessage
// @Failure 404 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /products [get]
func (p *productController) HandleGetProducts(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params, err := GetSearchParams(r)
	if err != nil {
		return err
	}

	products, err := p.productService.GetProducts(ctx, params)
	if err != nil {
		return domain.ConvertToWebErr(err)
	}

	return web.EncodeJSON(w, products, http.StatusOK)
}

// HandleGetProductByID godoc
// @Summary Get product by id
// @Description Get product by id
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path int true "product id"
// @Success 200 {object} dto.ProductDTO
// @Failure 400 {object} ErrorMessage
// @Failure 400 {object} ErrorMessage
// @Failure 404 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /product/{id} [get]
func (p *productController) HandleGetProductByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id, err := web.ParamInt(r, "id")
	if err != nil {
		return err
	}

	product, err := p.productService.FindById(ctx, int64(id))
	if err != nil {
		return domain.ConvertToWebErr(err)
	}

	return web.EncodeJSON(w, product, http.StatusOK)
}

// HandleCreateProduct godoc
// @Summary Create product
// @Description Create product
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body string true "product"
// @Success 201 {object} dto.ProductDTO
// @Failure 400 {object} ErrorMessage
// @Failure 404 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /product [post]
func (p *productController) HandleCreateProduct(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var productDTO dto.ProductDTO
	if err := json.NewDecoder(r.Body).Decode(&productDTO); err != nil {
		return errorhandling.NewBadRequestAPIError("invalid json body")
	}

	product, err := p.productService.CreateProduct(ctx, productDTO)
	if err != nil {
		return domain.ConvertToWebErr(err)
	}
	return web.EncodeJSON(w, product, http.StatusCreated)
}

// HandleFindProductByCategory godoc
// @Summary Find product by category
// @Description Find product by category
// @Tags products
// @Accept  json
// @Produce  json
// @Param category path string true "category"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} ErrorMessage
// @Failure 400 {object} ErrorMessage
// @Failure 404 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /products/category/{category} [get]
func (p *productController) HandleFindProductByCategory(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	category := web.Param(r, "category")

	c, err := p.productService.GetProductsByCategory(ctx, category)
	if err != nil {
		return domain.ConvertToWebErr(err)
	}

	return web.EncodeJSON(w, c, http.StatusOK)

}

// HandleUpdateProduct godoc
// @Summary Update product
// @Description Update product
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path int true "product id"
// @Param product body string true "product"
// @Success 200 {object} dto.ProductDTO
// @Failure 400 {object} ErrorMessage
// @Failure 400 {object} ErrorMessage
// @Failure 404 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /product/{id} [put]
func (p *productController) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var productDTO dto.ProductUpdateDTO

	productID, err := web.ParamInt(r, "id")
	if err != nil {
		apiErr := domain.NewBadRequest("invalid id", nil)
		log.Error(ctx, apiErr.Error())
		return domain.ConvertToWebErr(apiErr)
	}

	if err = json.NewDecoder(r.Body).Decode(&productDTO); err != nil {
		apiErr := domain.NewInternalError("error unmarshalling body", err)
		log.Error(ctx, apiErr.Error())
		return domain.ConvertToWebErr(apiErr)
	}

	product, err := p.productService.UpdateProduct(ctx, productDTO, int64(productID))
	if err != nil {
		return domain.ConvertToWebErr(err)
	}

	return web.EncodeJSON(w, product, http.StatusOK)
}

// HandleDeleteProduct godoc
// @Summary Delete product
// @Description Delete product
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path int true "product id"
// @Success 204
// @Failure 400 {object} ErrorMessage
// @Failure 400 {object} ErrorMessage
// @Failure 404 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /product/{id} [delete]
func (p *productController) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	productID, err := web.ParamInt(r, "id")
	if err != nil {
		apiErr := domain.NewBadRequest("invalid id", nil)
		log.Error(ctx, apiErr.Error())
		return domain.ConvertToWebErr(apiErr)
	}

	if err = p.productService.DeleteProduct(ctx, int64(productID)); err != nil {
		return domain.ConvertToWebErr(err)
	}

	return web.EncodeJSON(w, nil, http.StatusNoContent)
}

func GetSearchParams(r *http.Request) (dto.SearchParams, error) {
	params := dto.SearchParams{}
	query := r.URL.Query()
	if value := query.Get("sort"); value != "" {
		params.Sort = &value
	}
	if value := query.Get("title"); value != "" {
		params.Title = &value
	}

	if value := query.Get("name"); value != "" {
		params.Name = &value
	}

	if value := query.Get("limit"); value != "" {
		limit, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return params, errorhandling.NewRequestError(fmt.Sprintf("limit parameter value is not an integer. limit = %s", value))
		}
		params.Limit = &limit
	}
	if value := query.Get("offset"); value != "" {
		offset, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return params, errorhandling.NewRequestError(fmt.Sprintf("offset parameter value is not an integer. offset = %s", value))
		}
		params.Offset = &offset
	}

	if params.Limit == nil {
		defaultLimit := int64(config.LimitSearchRows)
		params.Limit = &defaultLimit
	}

	if params.Offset == nil {
		defaultOffset := int64(1)
		params.Offset = &defaultOffset
	}

	if value := query.Get("min"); value != "" {
		min, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return params, errorhandling.NewRequestError(fmt.Sprintf("min parameter value is not a float. min = %s", value))
		}
		params.Min = &min
	}
	if value := query.Get("max"); value != "" {
		max, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return params, errorhandling.NewRequestError(fmt.Sprintf("max parameter value is not a float. max = %s", value))
		}
		params.Max = &max
	}
	return params, nil
}
