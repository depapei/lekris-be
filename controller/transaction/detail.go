package transaction

import (
	helper "Lekris-BE/helpers"
	"Lekris-BE/model"
	res "Lekris-BE/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Detail(c *gin.Context) {
	id := c.Param("id")
	var trx model.Transaction

	if err := model.DB.Preload("DetailTransaction.Product").First(&trx, id).Error; err != nil {
		c.JSON(http.StatusNotFound, res.GeneralResponse{
			Success: false,
			Message: "Transaction not found",
		})
		return
	}

	c.JSON(http.StatusOK, res.GeneralResponse{
		Success: true,
		Data:    helper.MapTransactionToResponse(trx),
	})
}
