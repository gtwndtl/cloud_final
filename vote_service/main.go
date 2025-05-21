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

	api := r.Group("/api")
	{
		api.POST("/vote", controller.CreateVote)
		api.GET("/votes", controller.GetAllVotes)
		api.GET("/votes/candidate/:candidate_id", controller.GetVotesByCandidate)
		api.DELETE("/vote/:id", controller.DeleteVote)
	}

	r.Run(":8004") // ตัวอย่าง port 8005
}
