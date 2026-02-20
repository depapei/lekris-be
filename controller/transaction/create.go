package transaction

import (
	helper "Lekris-BE/helpers"
	"Lekris-BE/model"
	res "Lekris-BE/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	var input ValidateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// ... handle validation error (sama seperti kodemu) ...
		c.JSON(http.StatusBadRequest, gin.H{"Success": false, "Message": err.Error()})
		return
	}

	newTx := model.Transaction{
		Branchname:          input.Branchname,
		Totalprice:          input.Totalprice,
		Isreturningcustomer: input.Isreturningcustomer,
	}

	// Gunakan DB Transaction untuk atomicity
	dbTx := model.DB.Begin()
	if dbTx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Message": "Failed to begin transaction"})
		return
	}

	if err := dbTx.Create(&newTx).Error; err != nil {
		dbTx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Message": "Failed to create transaction", "Error": err.Error()})
		return
	}

	var details []model.DetailTransaction
	for _, item := range input.Items {
		details = append(details, model.DetailTransaction{
			TransactionID: newTx.ID,
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
		})
	}

	if len(details) > 0 {
		if err := dbTx.Create(&details).Error; err != nil {
			dbTx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Message": "Failed to create items", "Error": err.Error()})
			return
		}
	}

	if err := dbTx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Message": "Failed to commit"})
		return
	}

	// Fetch with preload for response
	var result model.Transaction
	model.DB.Preload("DetailTransaction.Product").First(&result, newTx.ID)

	c.JSON(http.StatusCreated, res.GeneralResponse{
		Success: true,
		Message: "Transaction created successfully 🎉",
		Data:    helper.MapTransactionToResponse(result),
	})
}
