package repositories

import (
	"sort"
	"time"

	"github.com/zranitm/sandbox/repositories/model"
)

type TransactionRepo interface {
	GetTxWithFilter(fromTime time.Time, toTime time.Time) ([]model.Transaction, int, int)
}

type transactionRepo struct {
	data []model.Transaction
}

func NewTransactionRepo(filePath string) (TransactionRepo, error) {
	data, err := model.GetAllTransactions(filePath)
	if err != nil {
		return nil, err
	}
	return &transactionRepo{
		data: data,
	}, nil
}

func (r *transactionRepo) GetTxWithFilter(fromTime time.Time, toTime time.Time) ([]model.Transaction, int, int) {
	var (
		result           []model.Transaction
		totalIncome      int
		totalExpenditure int
	)

	for _, tx := range r.data {
		if (tx.Date.Equal(fromTime) || tx.Date.After(fromTime)) && tx.Date.Before(toTime) {
			result = append(result, tx)
			if tx.Amount > 0 {
				totalIncome += tx.Amount
			} else {
				totalExpenditure += tx.Amount
			}
		}
	}

	// desc order by date
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.After(result[j].Date)
	})

	return result, totalIncome, totalExpenditure
}
