package main

import (
	"fmt"
	"time"

	"example.com/se/config"
	"example.com/se/controller"
	"example.com/se/metrics"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	PushGatewayURL = "http://pushgateway:9091" // เปลี่ยนเป็น URL ของ Pushgateway ของคุณ
	JobName        = "vote_service"
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		timer := prometheus.NewTimer(metrics.HTTPRequestDurationSeconds.WithLabelValues(method, path))
		defer timer.ObserveDuration()

		c.Next()

		status := c.Writer.Status()
		metrics.HTTPRequestsTotal.WithLabelValues(method, path, fmt.Sprintf("%d", status)).Inc()
	}
}

func pushMetricsLoop() {
	for {
		err := push.New(PushGatewayURL, JobName).
			Collector(metrics.HTTPRequestsTotal).
			Collector(metrics.HTTPRequestDurationSeconds).
			Push()
		if err != nil {
			fmt.Println("Failed to push metrics:", err)
		} else {
			fmt.Println("Metrics pushed to Pushgateway")
		}
		time.Sleep(10 * time.Second) // ปรับเวลาได้ตามต้องการ
	}
}

func main() {
	config.ConnectionDB()
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

	go pushMetricsLoop() // เรียกใช้งาน push metrics แบบ background goroutine

	r.Run(":8004")
}
