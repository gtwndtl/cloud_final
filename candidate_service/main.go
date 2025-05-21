package main

import (
	"log"
	"time"

	"example.com/se/config"
	"example.com/se/controller"
	"example.com/se/entity"
	"example.com/se/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	PushGatewayURL = "http://pushgateway:9091" // เปลี่ยนเป็น container name ถ้าใช้ Docker network
	JobName        = "candidate_service"
)

func main() {
	config.ConnectionDB()
	config.SetupDatabase()
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

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Loop push metrics ทุก 30 วินาที
	go func() {
		for {
			pushSnapshotMetrics()
			time.Sleep(30 * time.Second)
		}
	}()

	r.Run(":8003")
}

func pushSnapshotMetrics() {
	db := config.DB()

	var candidateCount int64
	db.Model(&entity.Candidates{}).Count(&candidateCount)
	metrics.CandidatesTotal.Set(float64(candidateCount))

	err := push.New(PushGatewayURL, JobName).
		Collector(metrics.CandidatesTotal).
		Push()

	if err != nil {
		log.Printf("❌ PushGateway error: %v", err)
	} else {
		log.Println("✅ Pushed candidate count to PushGateway")
	}
}
