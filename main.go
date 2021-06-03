package main

import (
	"gin-rest-api-example/config"
	"gin-rest-api-example/controllers"
	_ "gin-rest-api-example/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Gin Rest Api Example Swagger
// @version 1.0
// @description Gin Rest Api Example Swagger

// @contact.name Jeffrey Chu
// @contact.email jeffreychu888hk@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
func main() {
	r := gin.Default()

	// Connect to database
	config.ConnectDatabase()

	// Connect to redis
	config.ConnectRedis()

	// Routes
	// Book Router
	bookRouter := r.Group("")
	{
		bookRouter.GET("/books", controllers.FindBooks)
		bookRouter.GET("/books/:id", controllers.FindBook)
		bookRouter.POST("/books", controllers.CreateBook)
		bookRouter.PATCH("/books/:id", controllers.UpdateBook)
		bookRouter.DELETE("/books/:id", controllers.DeleteBook)
	}

	// Swagger
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Run the server
	r.Run()
}
