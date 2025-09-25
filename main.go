package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zranitm/sandbox/repositories"
	"github.com/zranitm/sandbox/services"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var timeInput time.Time
	for {
		fmt.Print("Enter history time (YYYYMM): ")
		historyTime, _ := reader.ReadString('\n')
		historyTime = strings.TrimSpace(historyTime)

		t, err := validateHistoryTime(historyTime)
		if err != nil {
			fmt.Println("Invalid history time:", err)
			continue
		}
		timeInput = t
		break
	}

	var filePath string
	for {
		fmt.Print("Enter file path: ")
		p, _ := reader.ReadString('\n')
		p = strings.TrimSpace(p)

		if err := validateFilePath(p, ".csv"); err != nil {
			fmt.Println("Invalid file path:", err)
			continue
		}
		filePath = p
		break
	}

	transactionRepo, err := repositories.NewTransactionRepo(filePath)
	if err != nil {
		fmt.Println("Error initializing transaction repository:", err)
		return
	}
	transactionService := services.NewTransactionService(transactionRepo)

	// if err := transactionService.ExportTxToCSV(timeInput, "./output/output.csv"); err != nil {
	// 	fmt.Println("Error exporting transactions:", err)
	// 	return
	// }
	// if err := transactionService.ExportTxToJSON(timeInput, "./output/output.json"); err != nil {
	// 	fmt.Println("Error exporting transactions:", err)
	// 	return
	// }

	resp := transactionService.GetTxWithFilter(timeInput)
	b, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}
	fmt.Println(string(b))
}

func validateHistoryTime(input string) (time.Time, error) {
	if len(input) != 6 {
		return time.Time{}, fmt.Errorf("invalid format: must be YYYYMM")
	}

	t, err := time.Parse("200601", input)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date: %w", err)
	}

	// since it's history time, it should not be in the future
	now := time.Now()
	if t.After(now) {
		return time.Time{}, fmt.Errorf("date cannot be in the future")
	}

	return t, nil
}

func validateFilePath(path string, expectedExt string) error {
	path = filepath.Clean(path)

	if expectedExt != "" && !strings.EqualFold(filepath.Ext(path), expectedExt) {
		return fmt.Errorf("file must have %s extension", expectedExt)
	}

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("cannot stat file %s: %w", path, err)
	}
	if fi.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open file for read: %w", err)
	}
	f.Close()

	return nil
}
