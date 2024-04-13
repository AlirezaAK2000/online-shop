package main

import (
	"github.com/AlirezaAK2000/online-shop/controllers"
	"github.com/AlirezaAK2000/online-shop/initializers"
	"github.com/AlirezaAK2000/online-shop/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	// initializers.EnvVariableInitializer()
	// initializers.InitializeMongoConnection()
	// initializers.InitilizeRedisConnection()
}

func main() {

	defer initializers.DisconnectMongoDBClient()
	r := gin.Default()
	r.POST("/product", middleware.RequireAuth, controllers.InsertProductController)
	r.GET("/product", middleware.RequireAuth, controllers.GetAllProductsController)
	r.GET("/product/:id", middleware.RequireAuth, controllers.GetProductByIDController)
	r.DELETE("/product/:id", middleware.RequireAuth, controllers.DeleteProductByIDController)
	r.POST("/signin", controllers.SignInController)
	r.POST("/login", controllers.LogInController)

	r.Run()
}
