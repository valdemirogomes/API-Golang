package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/domain"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/dto"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/service"
	"github.com/mercadolibre/fury_go-core/pkg/web"
)

type CategoryController interface {
	HandleCreateCategory(w http.ResponseWriter, r *http.Request) error
	HandleGetCategories(w http.ResponseWriter, r *http.Request) error
	HandleUpdateCategory(w http.ResponseWriter, r *http.Request) error
	HandleDeleteCategory(w http.ResponseWriter, r *http.Request) error
}

type categoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) CategoryController {
	return &categoryController{
		categoryService: categoryService,
	}
}

// HandleCreateCategory godoc
// @Summary Create category
// @Description Create category
// @Tags categories
// @Accept  json
// @Produce  json
// @Param category body string true "category"
// @Success 201 {object} domain.Category
// @Failure 400 {object} ErrorMessage
// @Failure 400 {object} ErrorMessage
// @Failure 404 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /category [post]
func (c *categoryController) HandleCreateCategory(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var categoryDomain domain.Category
	if err := json.NewDecoder(r.Body).Decode(&categoryDomain); err != nil {
		appErr := domain.NewInternalError("error unmarshalling request body", err)
		return domain.ConvertToWebErr(appErr)
	}
	category, err := c.categoryService.CreateCategory(ctx, categoryDomain)
	if err != nil {
		return domain.ConvertToWebErr(err)
	}

	return web.EncodeJSON(w, category, http.StatusCreated)
}

// HandleGetCategories godoc
// @Summary Get categories
// @Description Get categories
// @Tags categories
// @Accept  json
// @Produce  json
// @Param sort query string false "sort"
// @Param title query string false "title"
// @Param name query string false "name"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Param min query float64 false "min"
// @Param max query float64 false "max"
// @Success 200 {object} dto.CategoryListResponseDTO
// @Failure 400 {object} ErrorMessage
// @Failure 400 {object} ErrorMessage
// @Failure 404 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /products/categories [get]
func (c *categoryController) HandleGetCategories(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params, err := GetSearchParams(r)
	if err != nil {
		return err
	}

	categories, err := c.categoryService.GetCategories(ctx, params)
	if err != nil {
		return domain.ConvertToWebErr(err)
	}

	return web.EncodeJSON(w, categories, http.StatusOK)
}

// HandleUpdateCategory godoc
// @Summary Update category
// @Description Update category
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path int true "category id"
// @Param category body string true "category"
// @Success 200 {object} dto.CategoryDTO
// @Failure 400 {object} ErrorMessage
// @Failure 400 {object} ErrorMessage
// @Failure 404 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /category/{id} [put]
func (c *categoryController) HandleUpdateCategory(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var categoryDTO dto.CategoryDTO

	productID, err := web.ParamInt(r, "id")
	if err != nil {
		apiErr := domain.NewBadRequest(fmt.Sprintf("invalid product id: %d", productID), err)
		return domain.ConvertToWebErr(apiErr)
	}

	if err := json.NewDecoder(r.Body).Decode(&categoryDTO); err != nil {
		appErr := domain.NewInternalError("error unmarshalling request body", err)
		return domain.ConvertToWebErr(appErr)
	}
	category, err := c.categoryService.UpdateCategory(ctx, categoryDTO, int64(productID))
	if err != nil {
		return domain.ConvertToWebErr(err)
	}

	return web.EncodeJSON(w, category, http.StatusOK)
}

// HandleDeleteCategory godoc
// @Summary Delete category
// @Description Delete category
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path int true "category id"
// @Success 200
// @Failure 400 {object} ErrorMessage
// @Failure 400 {object} ErrorMessage
// @Failure 404 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /category/{id} [delete]
func (c *categoryController) HandleDeleteCategory(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	categoryID, err := web.ParamInt(r, "id")
	if err != nil {
		apiErr := domain.NewBadRequest(fmt.Sprintf("invalid category id: %d", categoryID), err)
		return domain.ConvertToWebErr(apiErr)
	}

	_, err = c.categoryService.DeleteCategory(ctx, int64(categoryID))
	if err != nil {
		return domain.ConvertToWebErr(err)
	}

	return web.EncodeJSON(w, nil, http.StatusOK)
}
