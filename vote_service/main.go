package main

import (
	"fmt"

	"example.com/se/config"
	"example.com/se/controller"
	"example.com/se/metrics"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		// ใช้ prometheus.NewTimer เพื่อจับเวลาการทำงาน
		timer := prometheus.NewTimer(metrics.HTTPRequestDurationSeconds.WithLabelValues(method, path))
		defer timer.ObserveDuration()

		c.Next()

		status := c.Writer.Status()
		metrics.HTTPRequestsTotal.WithLabelValues(method, path, fmt.Sprintf("%d", status)).Inc()
	}
}

func main() {
    config.ConnectionDB() // เรียกใช้งานตรง ๆ เพราะไม่มีค่า error คืนกลับมา
    config.SetupDatabase()

    metrics.RegisterMetrics()

    r := gin.Default()
    r.Use(cors.Default())
    r.Use(prometheusMiddleware())

    api := r.Group("/")
    {
        api.POST("/vote", controller.CreateVote)
        api.GET("/votes", controller.GetAllVotes)
        api.GET("/votes/candidate/:candidate_id", controller.GetVotesByCandidate)
        api.DELETE("/vote/:id", controller.DeleteVote)
    }

    r.GET("/metrics", gin.WrapH(promhttp.Handler()))

    r.Run(":8004")
}

