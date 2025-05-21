package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"example.com/se/config"
	"example.com/se/entity"
	"example.com/se/metrics"
	"example.com/se/services"
)

type SignUpPayload struct {
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Age       uint8     `json:"age" binding:"required"`
	Password  string    `json:"password" binding:"required,min=6"`
	Role      string    `json:"role"`
	BirthDay  time.Time `json:"birthday" binding:"required"`
	GenderID  uint      `json:"gender_id" binding:"required"`
}

type SignInPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func SignUp(c *gin.Context) {
	var payload SignUpPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		metrics.UserSignupFailuresTotal.Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()
	var existingUser entity.Users
	result := db.Where("email = ?", payload.Email).First(&existingUser)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		metrics.UserSignupFailuresTotal.Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if existingUser.ID != 0 {
		metrics.UserSignupFailuresTotal.Inc()
		c.JSON(http.StatusConflict, gin.H{"error": "Email is already registered"})
		return
	}

	hashedPassword, err := config.HashPassword(payload.Password)
	if err != nil {
		metrics.UserSignupFailuresTotal.Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := entity.Users{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Age:       payload.Age,
		Password:  hashedPassword,
		Role:      payload.Role,
		BirthDay:  payload.BirthDay,
		GenderID:  payload.GenderID,
	}

	if err := db.Create(&user).Error; err != nil {
		metrics.UserSignupFailuresTotal.Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metrics.UserSignupsTotal.Inc()

	var count int64
	db.Model(&entity.Users{}).Count(&count)
	metrics.UsersTotal.Set(float64(count))

	c.JSON(http.StatusCreated, gin.H{"message": "Sign-up successful"})
}

func SignIn(c *gin.Context) {
	var payload SignInPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		metrics.UserLoginFailuresTotal.Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user entity.Users
	err := config.DB().Where("email = ?", payload.Email).First(&user).Error
	if err != nil {
		metrics.UserLoginFailuresTotal.Inc()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		metrics.UserLoginFailuresTotal.Inc()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect"})
		return
	}

	metrics.UserLoginsTotal.Inc()

	jwtWrapper := services.JwtWrapper{
		SecretKey:       "SvNQpBN8y3qlVrsGAYYWoJJk56LtzFHx",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}

	token, err := jwtWrapper.GenerateToken(user.Email)
	if err != nil {
		metrics.UserLoginFailuresTotal.Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token_type": "Bearer",
		"token":      token,
		"id":         user.ID,
	})
}
