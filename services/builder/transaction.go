package builder

import (
	"time"

	intModel "github.com/zranitm/sandbox/repositories/model"
	extModel "github.com/zranitm/sandbox/services/model"
)

func BuildTransactionResponse(tx []intModel.Transaction, totalIncome int, totalExpenditure int, inputDate time.Time) extModel.TransactionResponse {
	var transactions []extModel.Transaction
	for _, t := range tx {
		transactions = append(transactions, extModel.Transaction{
			Date:    t.Date.Format("2006/01/02"),
			Amount:  t.Amount,
			Content: t.Content,
		})
	}

	return extModel.TransactionResponse{
		Period:           inputDate.Format("2006/01"),
		TotalIncome:      totalIncome,
		TotalExpenditure: totalExpenditure,
		Transactions:     transactions,
	}
}
