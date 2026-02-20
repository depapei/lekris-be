package product

import (
	"Lekris-BE/helpers"
	"Lekris-BE/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Update(c *gin.Context) {
	id := c.Param("id")

	var data model.Product

	if err := model.DB.First(&data, `"id" = ?`, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Success": false,
			"Message": "Record not found",
		})
		return
	}

	var input ValidateProductInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			// 422: Validation Errors
			out := make([]helpers.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = helpers.ErrorMsg{
					Field:   fe.Field(),
					Message: helpers.GetErrorMsg(fe),
				}
			}
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"Success": false,
				"Message": "Input validation failed",
				"Errors":  out,
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Success": false,
			"Message": "Invalid request body",
		})
		return
	}

	updateData := map[string]interface{}{
		"item":        input.Item,
		"description": input.Description,
		"price":       input.Price,
	}

	if err := model.DB.Model(&data).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Success": false,
			"Message": "Failed to update product due to database error.",
			"Error":   err.Error(),
		})
		return
	}
	response := data
	c.JSON(http.StatusOK, gin.H{
		"Success": true,
		"Message": "Product Updated Successfully",
		"Data":    response,
	})

}
