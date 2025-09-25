package services

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/zranitm/sandbox/repositories"
	"github.com/zranitm/sandbox/services/builder"
	extModel "github.com/zranitm/sandbox/services/model"
)

type TransactionService interface {
	GetTxWithFilter(inputDate time.Time) extModel.TransactionResponse
	ExportTxToCSV(inputDate time.Time, filePath string) error
	ExportTxToJSON(inputDate time.Time, filePath string) error
}

type transactionService struct {
	repo repositories.TransactionRepo
}

func NewTransactionService(repo repositories.TransactionRepo) TransactionService {
	return &transactionService{
		repo: repo,
	}
}

func (s *transactionService) GetTxWithFilter(inputDate time.Time) extModel.TransactionResponse {
	fromTime := inputDate
	toTime := inputDate.AddDate(0, 1, 0)

	transactions, totalIncome, totalExpenditure := s.repo.GetTxWithFilter(fromTime, toTime)

	return builder.BuildTransactionResponse(transactions, totalIncome, totalExpenditure, inputDate)
}

func (s *transactionService) ExportTxToCSV(inputDate time.Time, filePath string) error {
	resp := s.GetTxWithFilter(inputDate)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"date", "amount", "content"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, tx := range resp.Transactions {
		if err := writer.Write([]string{tx.Date, fmt.Sprintf("%d", tx.Amount), tx.Content}); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

func (s *transactionService) ExportTxToJSON(inputDate time.Time, filePath string) error {
	resp := s.GetTxWithFilter(inputDate)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(resp); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	return nil
}
