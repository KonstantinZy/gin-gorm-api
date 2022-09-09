package main

import (
	"github.com/KonstantinZy/gin-gorm-api/controllers"
	"github.com/KonstantinZy/gin-gorm-api/models"

	"github.com/gin-gonic/gin"
)

func main() {

	models.StartDB("api.db")
	// migrate all models to database
	models.MigrateOrm()

	r := gin.Default()
	controllers.RegisterTaskURI(r)

	r.Run()
}
