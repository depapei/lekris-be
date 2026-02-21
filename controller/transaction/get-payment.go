package transaction

import (
	"Lekris-BE/model"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPaymentProof(c *gin.Context) {
	// 1. Get transaction ID dari URL parameter
	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Success": false,
			"Message": "Transaction ID is required",
		})
		return
	}

	// 2. Query transaction dari database
	var transaction model.Transaction
	if err := model.DB.First(&transaction, transactionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"Success": false,
				"Message": "Transaction not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"Success": false,
			"Message": "Failed to fetch transaction",
			"Error":   err.Error(),
		})
		return
	}

	// 3. Cek apakah transaction memiliki payment proof
	if transaction.Imagepath == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"Success": false,
			"Message": "Payment proof not found for this transaction",
		})
		return
	}

	// 4. Convert relative path ke absolute path
	// Imagepath di DB: "/uploads/transactions/trx_abc123.png"
	// Absolute path: "./uploads/transactions/trx_abc123.png"
	absolutePath := filepath.Join(".", transaction.Imagepath)

	// 5. Cek apakah file benar-benar ada di filesystem
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"Success": false,
			"Message": "Payment proof file not found on server",
		})
		return
	}

	// 6. Serve file sebagai image
	// Opsi A: Inline (langsung tampil di browser)
	c.File(absolutePath)

	// Opsi B: Attachment (download) - uncomment jika mau download
	// filename := filepath.Base(transaction.Imagepath)
	// c.FileAttachment(absolutePath, filename)
}
