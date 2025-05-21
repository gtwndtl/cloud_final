package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/se/config"
	"example.com/se/controller"
)

const PORT = "8002"

func main() {
	config.ConnectionDB()
	config.SetupDatabase()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Election Service Running on PORT: %s", PORT)
	})

	api := r.Group("/api")
	{
		api.GET("/elections", controller.GetAllElections)
		api.GET("/election/:id", controller.GetElection)
		api.POST("/election", controller.CreateElection)
		api.PUT("/election/:id", controller.UpdateElection)
		api.DELETE("/election/:id", controller.DeleteElection)
	}

	r.Run(":" + PORT)
}
