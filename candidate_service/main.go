package main

import (
	"example.com/se/config"
	"example.com/se/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectionDB()
	config.SetupDatabase()

	r := gin.Default()

	api := r.Group("/")
	{
		api.GET("/candidates", controller.GetAllCandidates)
		api.GET("/candidate/:id", controller.GetCandidate)
		api.POST("/candidate", controller.CreateCandidate)
		api.PUT("/candidate/:id", controller.UpdateCandidate)
		api.DELETE("/candidate/:id", controller.DeleteCandidate)
	}

	r.Run(":8003")
}
