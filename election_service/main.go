package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"example.com/se/config"
	"example.com/se/controller"
	"example.com/se/metrics"
)

const PORT = "8002"


func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		observer := metrics.HTTPRequestDurationSeconds.WithLabelValues(method, path)
		timer := prometheus.NewTimer(observer)
		defer timer.ObserveDuration()

		c.Next()

		status := fmt.Sprintf("%d", c.Writer.Status())
		metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	}
}


func main() {
	config.ConnectionDB()
	config.SetupDatabase()

	metrics.RegisterMetrics()

	r := gin.Default()

	// ใช้ middleware Prometheus สำหรับจับ metrics
	r.Use(PrometheusMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Election Service Running on PORT: %s", PORT)
	})

	api := r.Group("/")
	{
		api.GET("/elections", controller.GetAllElections)
		api.GET("/election/:id", controller.GetElection)
		api.POST("/election", controller.CreateElection)
		api.PUT("/election/:id", controller.UpdateElection)
		api.DELETE("/election/:id", controller.DeleteElection)
	}

	// เพิ่ม endpoint สำหรับ Prometheus metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(":" + PORT)
}
