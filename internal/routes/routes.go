package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/handlers"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/middleware"
)

func SetRoutes(r *gin.Engine) {

	// API v1 prefix for all routes
	api := r.Group("/api/v1")

	//user routes

	user := api.Group("/users")
	{
		// ✅ Simple pattern: validation happens inside handler via ShouldBindJSON
		user.POST("/register", handlers.Register)
		user.POST("/login", handlers.Login)
		user.GET("/all", middleware.AuthMiddleware(), handlers.GetAllUsers)
		user.GET("/user/:id", middleware.AuthMiddleware(), handlers.GetUser)
		user.GET("/user", middleware.AuthMiddleware(), handlers.GetUserByEmail)
	}

	// product routes

	product := api.Group("/products")
	{
		product.GET("/all", handlers.GetAllProducts)

		productProtected := product.Group("/")
		productProtected.Use(middleware.AuthMiddleware())
		{
			productProtected.GET("/:id", handlers.GetProductById)
			productProtected.POST("/", handlers.CreateNewProduct)
			productProtected.PUT("/:id", handlers.UpdateProduct)
			productProtected.DELETE("/:id", handlers.DeleteProduct)
		}
	}

	// cart routes

	cart := api.Group(("/cart"))
	cartUserProtected := cart.Group("/")
	cartUserProtected.Use(middleware.AuthMiddleware())
	{
		cartUserProtected.GET("/items", handlers.GetAllCartItems)
		cartUserProtected.POST("/item", handlers.AddOrUpdateCartItem)
	}

	// Admin routes (admin authorized routes)
	cartAdminProtected := cart.Group("/")
	cartAdminProtected.Use(middleware.AuthMiddleware(), middleware.IsAuthorized("Admin"))
	{
		cartAdminProtected.GET("/:userId", handlers.GetCart)
		cartAdminProtected.DELETE("/:userId", handlers.DeleteCart)
	}
}
