package transaction

import (
	helper "Lekris-BE/helpers"
	"Lekris-BE/model"
	res "Lekris-BE/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Update(c *gin.Context) {
	id := c.Param("id")
	var input ValidateTransactionUpdate

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Success": false, "Message": err.Error()})
		return
	}

	var trx model.Transaction
	if err := model.DB.First(&trx, id).Error; err != nil {
		c.JSON(http.StatusNotFound, res.GeneralResponse{Success: false, Message: "Transaction not found"})
		return
	}

	dbTx := model.DB.Begin()
	if dbTx.Error != nil {
		c.JSON(http.StatusInternalServerError, res.GeneralResponse{Success: false, Message: "Failed to begin transaction"})
		return
	}

	// Update Header
	if err := dbTx.Model(&trx).Updates(map[string]interface{}{
		"Branchname":          input.Branchname,
		"Totalprice":          input.Totalprice,
		"Isreturningcustomer": input.Isreturningcustomer,
		"Customername":        input.Customername,
		"Updatedby":           input.UpdatedBy,
		"Updatedat":           time.Now(),
	}).Error; err != nil {
		dbTx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Message": "Failed to update header", "Error": err.Error()})
		return
	}

	// Delete old details, then insert new ones
	if err := dbTx.Where("transaction_id = ?", trx.ID).Delete(&model.DetailTransaction{}).Error; err != nil {
		dbTx.Rollback()
		c.JSON(http.StatusInternalServerError, res.GeneralResponse{Success: false, Message: "Failed to clear old items"})
		return
	}

	if len(input.Items) > 0 {
		var newDetails []model.DetailTransaction
		for _, item := range input.Items {
			newDetails = append(newDetails, model.DetailTransaction{
				TransactionID: trx.ID,
				ProductID:     item.ProductID,
				Quantity:      item.Quantity,
			})
		}
		if err := dbTx.Create(&newDetails).Error; err != nil {
			dbTx.Rollback()
			c.JSON(http.StatusInternalServerError, res.GeneralResponse{Success: false, Message: "Failed to insert new items"})
			return
		}
	}

	if err := dbTx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, res.GeneralResponse{Success: false, Message: "Failed to commit"})
		return
	}

	var updated model.Transaction
	model.DB.Preload("DetailTransaction.Product").First(&updated, trx.ID)

	c.JSON(http.StatusOK, res.GeneralResponse{
		Success: true,
		Message: "Transaction updated successfully",
		Data:    helper.MapTransactionToResponse(updated),
	})
}
