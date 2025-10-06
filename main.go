package main

import (
	"crud-api-golang/controllers"
	"crud-api-golang/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()
	router.POST("/posts", controllers.PostsCreate)
	router.PUT("/posts/:id", controllers.PostsUpdate)
	router.GET("/posts/:id", controllers.PostsShow)
	router.GET("/posts", controllers.PostsIndex)
	router.DELETE("/posts/:id", controllers.PostsDelete)

	router.Run()
}
