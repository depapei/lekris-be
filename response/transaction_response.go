package response

import "time"

type ItemResponse struct {
	ID          int32   `json:"id"`
	Item        string  `json:"item"`        // Dari Product.Item
	Description string  `json:"description"` // Dari Product.Description
	Price       int64   `json:"price"`       // Dari Product.Price
	Quantity    float64 `json:"quantity"`    // Dari DetailTransaction.Quantity
}

type TransactionResponse struct {
	ID                  int32          `json:"id"`
	Branchname          string         `json:"branchname"`
	Timestamp           time.Time      `json:"timestamp"`
	Totalprice          int64          `json:"totalprice"`
	Isreturningcustomer *bool          `json:"isreturningcustomer"`
	Items               []ItemResponse `json:"items"`
}

type GeneralResponse struct {
	Data    interface{} `json:"Data"` // Bisa []TransactionResponse atau TransactionResponse
	Success bool        `json:"Success"`
	Message string      `json:"Message,omitempty"` // Optional, hanya untuk error/info
}
