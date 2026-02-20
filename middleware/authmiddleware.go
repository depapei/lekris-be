package middleware

import (
	"Lekris-BE/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// fmt.Println("PATH:", c.FullPath())
		// fmt.Println("METHOD:", c.Request.Method)
		// fmt.Println("AUTH:", c.GetHeader("Authorization"))

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Please login first!"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Wrong token format!"})
			return
		}

		claims, err := utils.ParseJWT(parts[1])
		if err != nil {
			fmt.Println(utils.ParseJWT(parts[1]))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token!"})
			return
		}

		// simpan data user di context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
