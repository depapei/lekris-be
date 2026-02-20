package supply

import (
	"Lekris-BE/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidateSupplyInput struct {
	Name string `json:"name" binding:"required" example:"Minyak"`
	Unit string `json:"unit" binding:"required" example:"liter"`
}

func Index(c *gin.Context) {
	var data []model.Supply
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
