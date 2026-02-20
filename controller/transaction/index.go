package transaction

import (
	helper "Lekris-BE/helpers"
	"Lekris-BE/model"
	res "Lekris-BE/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionItemInput struct {
	ProductID int32   `json:"product_id" binding:"required" example:"1"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0" example:"2"`
}

type ValidateTransactionInput struct {
	Branchname          string                 `json:"branchName" binding:"required" example:"Pasar Segar"`
	Totalprice          int64                  `json:"totalPrice" binding:"required" example:"16000"`
	Isreturningcustomer *bool                  `json:"Isreturningcustomer" binding:"required"`
	Items               []TransactionItemInput `json:"items" binding:"required,dive"`
}

func Index(c *gin.Context) {
	var transactions []model.Transaction

	result := model.DB.Order("timestamp desc").Preload("DetailTransaction.Product").Limit(50).Find(&transactions)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, res.GeneralResponse{
			Success: false,
			Message: result.Error.Error(),
		})
		return
	}

	// Mapping ke Response DTO
	var responseList []res.TransactionResponse
	for _, trx := range transactions {
		responseList = append(responseList, helper.MapTransactionToResponse(trx))
	}

	c.JSON(http.StatusOK, res.GeneralResponse{
		Success: true,
		Data:    responseList, // Akan jadi [] jika kosong, konsisten!
	})
}
