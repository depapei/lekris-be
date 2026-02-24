package transaction

import (
	"Lekris-BE/helpers"
	helper "Lekris-BE/helpers"
	"Lekris-BE/model"
	res "Lekris-BE/response"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	var input ValidateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Success": false, "Message": err.Error()})
		return
	}

	// 1. Handle payment_proof upload DULUAN (sebelum DB transaction)
	now := time.Now().Format("02-01-2006")
	uploadDir := filepath.Join("uploads", "transactions", now)
	baseURL := "/uploads/transactions/" + now

	uploadResult, err := helpers.SaveBase64Image(input.PaymentProof, uploadDir, baseURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Success": false,
			"Message": fmt.Sprintf("Failed to process payment proof: %v", err),
		})
		return
	}
	// uploadResult.AbsolutePath akan kita gunakan untuk cleanup jika DB transaction gagal

	// 2. Siapkan data transaksi
	newTx := model.Transaction{
		Branchname:          input.Branchname,
		Totalprice:          input.Totalprice,
		Isreturningcustomer: input.Isreturningcustomer,
		Createdby:           int32(input.CreatedBy),
		Updatedby:           int32(input.CreatedBy),
		Customername:        input.Customername,
		Imagepath:           uploadResult.RelativePath, // Sudah dapat path dari upload
	}

	// 3. Mulai database transaction
	dbTx := model.DB.Begin()
	if dbTx.Error != nil {
		// Cleanup file karena DB transaction tidak bisa dimulai
		helpers.CleanupUploadedFile(uploadResult.AbsolutePath)
		c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Message": "Failed to begin transaction"})
		return
	}

	// 4. Insert header transaction
	if err := dbTx.Create(&newTx).Error; err != nil {
		dbTx.Rollback()
		helpers.CleanupUploadedFile(uploadResult.AbsolutePath) // Cleanup file jika DB gagal
		c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Message": "Failed to create transaction", "Error": err.Error()})
		return
	}

	// 5. Insert detail transaction
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
			helpers.CleanupUploadedFile(uploadResult.AbsolutePath) // Cleanup file jika DB gagal
			c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Message": "Failed to create items", "Error": err.Error()})
			return
		}
	}

	// 6. Commit semua perubahan
	if err := dbTx.Commit().Error; err != nil {
		helpers.CleanupUploadedFile(uploadResult.AbsolutePath) // Cleanup file jika commit gagal
		c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Message": "Failed to commit transaction"})
		return
	}

	// 7. Fetch result with preload untuk response
	var result model.Transaction
	model.DB.Preload("DetailTransaction.Product").First(&result, newTx.ID)

	c.JSON(http.StatusCreated, res.GeneralResponse{
		Success: true,
		Message: "Transaction created successfully 🎉",
		Data:    helper.MapTransactionToResponse(result),
	})
}
