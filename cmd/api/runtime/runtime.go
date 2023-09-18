package runtime

import (
	"context"

	"github.com/mercadolibre/fury_go-core/pkg/log"

	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/config"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/controller"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/lib/mysql"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/repository"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/service"
)

type Runtime struct {
	Environment        config.Environment
	ProductController  controller.ProductController
	CategoryController controller.CategoryController
}

func InstanceRuntime() *Runtime {
	env := config.InitConfig()

	//database
	mySQLClient, err := mysql.NewMySQL(env)
	if err != nil {
		log.Panic(context.Background(), err.Error())
	}

	//repositories
	productRepository := repository.NewProductRepository()
	categoryRepository := repository.NewCategoryRepository()

	//services
	productService := service.NewProductService(productRepository, categoryRepository, mySQLClient, env)
	categoryService := service.NewCategoryService(categoryRepository, mySQLClient, env)

	//controllers
	productController := controller.NewProductController(productService, env)
	categoryController := controller.NewCategoryController(categoryService)

	return &Runtime{
		Environment:        env,
		ProductController:  productController,
		CategoryController: categoryController,
	}
}
