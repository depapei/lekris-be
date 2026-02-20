package product

import (
	"Lekris-BE/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidateProductInput struct {
	Item        string `json:"item" binding:"required" example:"Lele Goreng"`
	Price       int    `json:"price" binding:"required"`
	Description string `json:"description"`
}

func Index(c *gin.Context) {
	var data []model.Product

	result := model.DB.Find(&data).Limit(100)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": result.Error.Error(),
		})
		return
	}

	if len(data) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"Success": true,
			"Message": "No Data Found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Success": true,
		"Data":    data,
	})
}
