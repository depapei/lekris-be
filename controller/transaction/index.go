package transaction

import (
	helper "Lekris-BE/helpers"
	"Lekris-BE/model"
	res "Lekris-BE/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionItemInput struct {
	ProductID int32   `json:"product_id" binding:"required" example:"1"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0" example:"2"`
}

type ValidateTransactionInput struct {
	Branchname          string                 `json:"branchName" binding:"required" example:"Pasar Segar"`
	Customername        string                 `json:"customerName" binding:"omitempty" example:"Rangga"`
	Totalprice          int64                  `json:"totalPrice" binding:"required" example:"16000"`
	Isreturningcustomer *bool                  `json:"isReturningCustomer" binding:"required"`
	CreatedBy           int64                  `json:"createdBy" binding:"required"`
	PaymentProof        string                 `json:"payment_proof" binding:"omitempty"`
	Items               []TransactionItemInput `json:"items" binding:"required,dive"`
}

type ValidateTransactionUpdate struct {
	Branchname          string                 `json:"branchName" binding:"required" example:"Pasar Segar"`
	Customername        string                 `json:"customerName" binding:"omitempty" example:"Rangga"`
	Totalprice          int64                  `json:"totalPrice" binding:"required" example:"16000"`
	Isreturningcustomer *bool                  `json:"isReturningCustomer" binding:"required"`
	UpdatedBy           int64                  `json:"updatedBy" binding:"required"`
	PaymentProof        string                 `json:"payment_proof" binding:"omitempty"`
	Items               []TransactionItemInput `json:"items" binding:"required,dive"`
}

func Index(c *gin.Context) {
	var transactions []model.Transaction

	result := model.DB.Where("isdelete = ?", time.Time{}).Order("timestamp desc").Preload("DetailTransaction.Product").Preload("CreatedByUser").Limit(50).Find(&transactions)

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
