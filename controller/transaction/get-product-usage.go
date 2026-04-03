package transaction

import (
	"Lekris-BE/model"
	res "Lekris-BE/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductUsageResult struct {
	Item  string `json:"item"`
	Total int64  `json:"total"`
}

func ProductUsage(c *gin.Context) {
	var results []ProductUsageResult
	date := c.Query("date")

	startDate, _ := time.Parse("2006-01-02", date)
	endDate := startDate.Add(24 * time.Hour)

	result := model.DB.
		Table("transactions t").
		Select("p.item, COUNT(*) as total").
		Joins("JOIN detail_transaction dt ON dt.transaction_id = t.id").
		Joins("JOIN products p ON p.id = dt.product_id").
		Where("t.timestamp >= ? AND t.timestamp < ?", startDate, endDate).
		Group("p.id, p.item").
		Scan(&results)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, res.GeneralResponse{
			Success: false,
			Message: result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.GeneralResponse{
		Success: true,
		Data:    results,
	})
}
