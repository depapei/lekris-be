package helpers

import (
	"Lekris-BE/model"
	"Lekris-BE/response"
)

func MapTransactionToResponse(trx model.Transaction) response.TransactionResponse {
	var items []response.ItemResponse
	for _, detail := range trx.DetailTransaction {
		items = append(items, response.ItemResponse{
			ID:          detail.ProductID,
			Item:        detail.Product.Item,
			Description: detail.Product.Description,
			Price:       detail.Product.Price,
			Quantity:    detail.Quantity,
		})
	}

	return response.TransactionResponse{
		ID:                  trx.ID,
		Branchname:          trx.Branchname,
		Timestamp:           trx.Timestamp,
		Totalprice:          trx.Totalprice,
		Isreturningcustomer: trx.Isreturningcustomer,
		Items:               items,
	}
}
