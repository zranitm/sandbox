package model

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Date    time.Time
	Amount  int
	Content string
}

func GetAllTransactions(filePath string) ([]Transaction, error) {
	tx, err := loadTransactionsCSV(filePath)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func loadTransactionsCSV(path string) ([]Transaction, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer f.Close()
	return readTransactionsCSV(bufio.NewReader(f))
}

func readTransactionsCSV(r io.Reader) ([]Transaction, error) {
	cr := csv.NewReader(r)

	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	want := []string{"date", "amount", "content"}
	if len(header) < 3 ||
		!strings.EqualFold(header[0], want[0]) ||
		!strings.EqualFold(header[1], want[1]) ||
		!strings.EqualFold(header[2], want[2]) {
		return nil, fmt.Errorf("unexpected header: %v (want %v)", header, want)
	}

	var out []Transaction
	line := 2 // data starts on line 2
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		if len(rec) == 0 {
			line++
			continue
		}
		if len(rec) < 3 {
			return nil, fmt.Errorf("line %d: expected 3 fields, got %d", line, len(rec))
		}

		dateStr := strings.TrimSpace(rec[0])
		amtStr := strings.TrimSpace(rec[1])
		contentStr := strings.TrimSpace(rec[2])

		d, err := time.Parse("2006/01/02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid date %q: %w", line, dateStr, err)
		}
		amt, err := strconv.Atoi(amtStr)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid amount %q: %w", line, amtStr, err)
		}

		out = append(out, Transaction{
			Date:    d,
			Amount:  amt,
			Content: contentStr,
		})
		line++
	}
	return out, nil
}
