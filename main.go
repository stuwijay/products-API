package main

import (
	"warungjwt_postgre/config"
	"warungjwt_postgre/handlers"
	"warungjwt_postgre/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.InitDatabase()

	e := echo.New()

	// Middleware logger
	e.Use(middleware.Logger())

	// Middleware recovery (optional, handles panics)
	e.Use(middleware.Recover())

	// Public routes
	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)

	// Protected routes - Role-based access
	products := e.Group("/products")
	products.Use(middlewares.IsAuthorized("admin")) // Middleware to allow only admin for POST, PUT, DELETE
	products.GET("", handlers.GetProducts)          // Accessible by both admin and staff
	products.GET("/:id", handlers.GetProduct)       // Accessible by both admin and staff
	products.POST("", handlers.CreateProduct)
	products.PUT("/:id", handlers.UpdateProduct)
	products.DELETE("/:id", handlers.DeleteProduct)

	// Route for staff (only GET access)
	staffProducts := e.Group("/staff/products")
	staffProducts.Use(middlewares.IsAuthorized("staff")) // Middleware to allow only staff with GET
	staffProducts.GET("", handlers.GetProducts)
	staffProducts.GET("/:id", handlers.GetProduct)

	e.Logger.Fatal(e.Start(":1323"))
}
