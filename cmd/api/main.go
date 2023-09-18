package main

import (
	"log"
	"time"

	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/config"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/runtime"

	"github.com/mercadolibre/fury_go-core/pkg/web"
	"github.com/mercadolibre/fury_go-platform/pkg/fury"
)

// @title           Swagger GO Ready Bases API
// @version         1.0
// @description     Project to reprocess flows.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
const (
	timeOutGeneral = 60
)

func routes(app *fury.Application, run *runtime.Runtime) {
	app.Use(config.JSONResponse())

	//Product
	app.Get("/products", run.ProductController.HandleGetProducts)
	app.Get("/product/{id}", run.ProductController.HandleGetProductByID)
	app.Get("/products/category/{category}", run.ProductController.HandleFindProductByCategory)
	app.Post("/product", run.ProductController.HandleCreateProduct)
	app.Put("/product/{id}", run.ProductController.HandleUpdateProduct)
	app.Delete("/product/{id}", run.ProductController.HandleDeleteProduct)

	//Category
	app.Get("/products/categories", run.CategoryController.HandleGetCategories)
	app.Post("/category", run.CategoryController.HandleCreateCategory)
	app.Put("/category/{id}", run.CategoryController.HandleUpdateCategory)
	app.Delete("/category/{id}", run.CategoryController.HandleDeleteCategory)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	run := runtime.InstanceRuntime()
	app, err := fury.NewWebApplication(
		fury.WithTimeouts(
			web.Timeouts{
				WriteTimeout:      timeOutGeneral * time.Minute,
				ReadTimeout:       timeOutGeneral * time.Minute,
				ReadHeaderTimeout: timeOutGeneral * time.Minute,
				IdleTimeout:       timeOutGeneral * time.Minute,
				ShutdownTimeout:   timeOutGeneral * time.Minute,
			},
		),
		fury.WithLogLevel(run.Environment.LogLevel),
	)

	if err != nil {
		return err
	}

	routes(app, run)
	return app.Run()

}
