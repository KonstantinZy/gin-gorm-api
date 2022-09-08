package main

import (
	"example/api/controllers"
	"example/api/models"

	"github.com/gin-gonic/gin"
)

func main() {

	// migrate all models to database
	models.MigrateOrm()

	r := gin.Default()

	r.POST("/task", controllers.CreateOneTask)
	r.GET("/task", controllers.Get)
	r.GET("/task/:id", controllers.GetOneTask)
	r.PATCH("/task/:id", controllers.UpdateTask)
	r.DELETE("/task/:id", controllers.DeleteTask)

	r.Run()
}
