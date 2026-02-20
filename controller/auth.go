package controller

import (
	"Lekris-BE/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type ValidateUserInput struct {
	Username string `json:"username" binding:"required" example:"admin_dev"`
	Password string `json:"password" binding:"required" example:"victorainitialproject"`
}

type JWTClaim struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	Sub      string `json:"sub"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

func Login(c *gin.Context) {
	var input ValidateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// cek user di DB
	var user model.User
	if err := model.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// validasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// bikin JWT token
	expirationTime := time.Now().Add(24 * time.Hour) // expired 24 jam
	claims := &JWTClaim{
		UserID:   user.ID,
		Sub:      (user.Username + "_" + string(user.ID)),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// Issuer:    "yourapp", // opsional
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Success":  true,
		"Message":  "Login success",
		"Token":    tokenString,
		"Username": user.Username,
		"Id":       user.ID,
	})
}
