package main

import (
	"example.com/se/config"
	"example.com/se/controller"
	"example.com/se/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config.ConnectionDB()
	config.SetupDatabase()

	// ลงทะเบียน metrics
	metrics.RegisterMetrics()

	r := gin.Default()

	api := r.Group("/")
	{
		api.GET("/candidates", controller.GetAllCandidates)
		api.GET("/candidate/:id", controller.GetCandidate)
		api.POST("/candidate", controller.CreateCandidate)
		api.PUT("/candidate/:id", controller.UpdateCandidate)
		api.DELETE("/candidate/:id", controller.DeleteCandidate)
	}

	// เพิ่ม endpoint สำหรับ metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(":8003")
}
