package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type transferPayload struct {
	IdempotencyKey string `json:"idempotency_key"`
	BankCode       string `json:"bank_code"`
	AccountNumber  string `json:"account_number"`
	Amount         int64  `json:"amount"`
}

func main() {
	targetURL := flag.String("url", "http://localhost:8080/transfer", "Target transfer endpoint URL")
	rps := flag.Int("rps", 20, "Target requests per second (minimum 10)")
	duration := flag.Duration("duration", 5*time.Second, "Test duration")
	amount := flag.Int64("amount", 100000, "Transfer amount per request")
	flag.Parse()

	if *rps < 1 {
		*rps = 1
	}

	totalRequests := int(duration.Seconds()) * (*rps)
	if totalRequests < 1 {
		totalRequests = *rps
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}

	var (
		acceptedCount     atomic.Int64
		insufficientCount atomic.Int64
		otherErrorCount   atomic.Int64
		networkErrorCount atomic.Int64
		totalLatencyNs    atomic.Int64
		maxLatencyNs      atomic.Int64
	)

	fmt.Printf("Starting load test: url=%s target_rps=%d duration=%s total_requests=%d amount=%d\n",
		*targetURL, *rps, *duration, totalRequests, *amount)

	interval := time.Second / time.Duration(*rps)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < totalRequests; i++ {
		<-ticker.C
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			reqBody := transferPayload{
				IdempotencyKey: fmt.Sprintf("loadtest-%s-%d", uuid.NewString(), idx),
				BankCode:       "BCA",
				AccountNumber:  "1234567890",
				Amount:         *amount,
			}
			bodyBytes, _ := json.Marshal(reqBody)

			reqStart := time.Now()
			resp, err := client.Post(*targetURL, "application/json", bytes.NewReader(bodyBytes))
			latNs := time.Since(reqStart).Nanoseconds()
			totalLatencyNs.Add(latNs)

			for {
				curMax := maxLatencyNs.Load()
				if latNs <= curMax || maxLatencyNs.CompareAndSwap(curMax, latNs) {
					break
				}
			}

			if err != nil {
				networkErrorCount.Add(1)
				return
			}
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusAccepted:
				acceptedCount.Add(1)
			case http.StatusUnprocessableEntity:
				insufficientCount.Add(1)
			default:
				otherErrorCount.Add(1)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)
	actualRPS := float64(totalRequests) / elapsed.Seconds()
	avgLatency := time.Duration(totalLatencyNs.Load() / int64(totalRequests))
	maxLatency := time.Duration(maxLatencyNs.Load())

	fmt.Println("\n--- Load Test Summary ---")
	fmt.Printf("Elapsed          : %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("Total Requests   : %d\n", totalRequests)
	fmt.Printf("Actual RPS       : %.2f req/s\n", actualRPS)
	fmt.Printf("202 Accepted     : %d\n", acceptedCount.Load())
	fmt.Printf("422 Insufficient : %d\n", insufficientCount.Load())
	fmt.Printf("Other HTTP Err   : %d\n", otherErrorCount.Load())
	fmt.Printf("Network Err      : %d\n", networkErrorCount.Load())
	fmt.Printf("Avg Latency      : %s\n", avgLatency.Round(time.Millisecond))
	fmt.Printf("Max Latency      : %s\n", maxLatency.Round(time.Millisecond))
}
