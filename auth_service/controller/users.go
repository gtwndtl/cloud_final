package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"example.com/se/config"
	"example.com/se/entity"
	"example.com/se/services"
)

const pushGatewayURL = "http://pushgateway:9091"
const jobName = "auth_service"

func pushMetric(collector prometheus.Collector) {
	if err := push.New(pushGatewayURL, jobName).
		Collector(collector).
		Push(); err != nil {
		log.Println("Failed to push metric:", err)
	}
}

func GetAllUsers(c *gin.Context) {
	db := config.DB()

	var users []entity.Users
	if err := db.Preload("Gender").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	var user entity.Users
	if err := db.Preload("Gender").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type UpdateUserPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       uint8  `json:"age"`
	Role      string `json:"role"`
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var payload UpdateUserPayload
	db := config.DB()

	// Check if user exists
	var user entity.Users
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Bind JSON payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Update fields
	user.FirstName = payload.FirstName
	user.LastName = payload.LastName
	user.Age = payload.Age
	user.Role = payload.Role

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Push metric for update
	updateGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "users_update_total",
		Help: "Total number of user updates",
	})
	updateGauge.Inc()
	pushMetric(updateGauge)

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	result := db.Delete(&entity.Users{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Push metric for delete
	deleteCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "users_delete_total",
		Help: "Total number of user deletes",
	})
	deleteCounter.Inc()
	pushMetric(deleteCounter)

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, BEARER_SCHEMA) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			return
		}

		tokenString := authHeader[len(BEARER_SCHEMA):]

		jwtWrapper := services.JwtWrapper{
			SecretKey:       "SvNQpBN8y3qlVrsGAYYWoJJk56LtzFHx",
			Issuer:          "AuthService",
			ExpirationHours: 24,
		}

		claims, err := jwtWrapper.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Save claims info in context for later use
		c.Set("email", claims.Email)

		c.Next()
	}
}