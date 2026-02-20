package supply

import (
	"Lekris-BE/helpers"
	"Lekris-BE/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Create(c *gin.Context) {
	var input ValidateSupplyInput

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
	data := model.Supply{
		Name: input.Name,
		Unit: input.Unit,
	}

	if err := model.DB.Create(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Success": false,
			"Message": "Failed to create supply due to database error.",
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"Success": true,
		"Message": "Supply data master created successfully. 🎉",
		"Data":    data,
	})
}
