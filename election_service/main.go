package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"

	"example.com/se/config"
	"example.com/se/controller"
	"example.com/se/metrics"
)

const (
	PORT          = "8002"
	PushGatewayURL = "http://pushgateway:9091"  // เปลี่ยนเป็น URL ของ Pushgateway ของคุณ
	JobName        = "election_service"       // ชื่อ job ที่จะส่ง metrics ไป
)

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

func pushMetrics() {
	// ส่ง metrics ไปที่ Pushgateway ทุก ๆ 10 วินาที (หรือจะปรับตามที่ต้องการ)
	for {
		err := push.New(PushGatewayURL, JobName).
			Collector(metrics.HTTPRequestsTotal).
			Collector(metrics.HTTPRequestDurationSeconds).
			Push()
		if err != nil {
			fmt.Println("Could not push metrics to Pushgateway:", err)
		} else {
			fmt.Println("Pushed metrics to Pushgateway")
		}
		time.Sleep(10 * time.Second)
	}
}

func main() {
	config.ConnectionDB()
	config.SetupDatabase()

	metrics.RegisterMetrics()

	go pushMetrics()  // รัน push metrics เป็น goroutine background

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

	// ยังคงมี endpoint /metrics ให้ scrape ได้ตามปกติ
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(":" + PORT)
}
