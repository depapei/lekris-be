package transaction

import (
	"Lekris-BE/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	id := c.Param("id")

	var data model.Transaction

	if err := model.DB.First(&data, `"id" = ?`, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{ // FIX: 404
			"Success": false,
			"Message": "Record not found",
		})
		return
	}

	if err := model.DB.
		Model(&model.Transaction{}).
		Delete(`"id" = ? `, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Success": false,
			"Message": "Failed to delete transaction due to database error.",
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Success": true,
		"Message": "Transaction Deleted Successfully",
	})
}
