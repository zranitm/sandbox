package model

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
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

	var (
		out        []Transaction
		line       = 2
		numWorkers = 5
		jobs       = make(chan []string, 10)
		results    = make(chan Transaction, 10)
		wg, ctx    = errgroup.WithContext(context.Background())
	)

	worker := func() error {
		for {
			select {
			case rec, ok := <-jobs:
				if !ok {
					return nil
				}
				dateStr := strings.TrimSpace(rec[0])
				amtStr := strings.TrimSpace(rec[1])
				contentStr := strings.TrimSpace(rec[2])

				d, err := time.Parse("2006/01/02", dateStr)
				if err != nil {
					return fmt.Errorf("invalid date %q: %w", dateStr, err)
				}
				amt, err := strconv.Atoi(amtStr)
				if err != nil {
					return fmt.Errorf("invalid amount %q: %w", amtStr, err)
				}

				results <- Transaction{
					Date:    d,
					Amount:  amt,
					Content: contentStr,
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	for i := 0; i < numWorkers; i++ {
		wg.Go(worker)
	}

	// collect results on the fly
	wgResult := sync.WaitGroup{}
	wgResult.Add(1)
	go func() {
		defer wgResult.Done()
		for {
			select {
			case tx, ok := <-results:
				if !ok {
					return
				}
				out = append(out, tx)
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Go(func() error {
		for {
			rec, err := cr.Read()
			if err == io.EOF {
				close(jobs)
				break
			}
			if err != nil {
				return fmt.Errorf("line %d: %w", line, err)
			}
			if len(rec) == 0 {
				line++
				continue
			}
			if len(rec) < 3 {
				return fmt.Errorf("line %d: expected 3 fields, got %d", line, len(rec))
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				jobs <- rec
			}
			line++
		}
		return nil
	})

	if err := wg.Wait(); err != nil {
		return nil, err
	}
	close(results)

	wgResult.Wait()
	return out, nil
}
