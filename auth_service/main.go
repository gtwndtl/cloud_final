package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"

	"example.com/se/config"
	"example.com/se/controller"
	"example.com/se/entity"
	"example.com/se/metrics"
)

const (
	PORT           = "8001"
	PushGatewayURL = "http://pushgateway:9091"
	JobName        = "auth_service"
)

func main() {
	config.ConnectionDB()
	config.SetupDatabase()
	metrics.RegisterMetrics()

	r := gin.Default()

	r.Use(controller.CORSMiddleware())
	
	r.Use(func(c *gin.Context) {
		timer := prometheus.NewTimer(metrics.HTTPRequestDurationSeconds.WithLabelValues(c.Request.Method, c.FullPath()))
		defer timer.ObserveDuration()

		c.Next()

		metrics.HTTPRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), http.StatusText(c.Writer.Status())).Inc()
	})



	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := r.Group("/api")
	{
		api.POST("/signup", controller.SignUp)
		api.POST("/signin", controller.SignIn)

		api.Use(controller.AuthorizeJWT())
		{
			api.GET("/users", controller.GetAllUsers)
			api.GET("/user/:id", controller.GetUser)
			api.PUT("/user/:id", controller.UpdateUser)
			api.DELETE("/user/:id", controller.DeleteUser)
		}
	}

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Auth Service Running on PORT: %s", PORT)
	})

	go func() {
		for {
			pushSnapshotMetrics()
			time.Sleep(30 * time.Second)
		}
	}()

	r.Run(":" + PORT)
}

func pushSnapshotMetrics() {
	db := config.DB()
	var userCount int64
	db.Model(&entity.Users{}).Count(&userCount)
	metrics.UsersTotal.Set(float64(userCount))

	err := push.New(PushGatewayURL, JobName).
		Collector(metrics.UsersTotal).
		Push()
	if err != nil {
		// log error silently
	}
}
