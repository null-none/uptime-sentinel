package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type CheckTask struct {
	ID  int
	URL string
}

type CheckResult struct {
	TaskID     int
	URL        string
	OK         bool
	StatusCode int
	Err        error
	Duration   time.Duration
}

const workerCount = 5
const requestTimeout = 5 * time.Second

func worker(_ int, jobs <-chan CheckTask, results chan<- CheckResult, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}

	for task := range jobs {
		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, task.URL, nil)
		if err != nil {
			cancel()
			results <- CheckResult{
				TaskID: task.ID,
				URL:    task.URL,
				OK:     false,
				Err:    fmt.Errorf("create request: %w", err),
			}
			continue
		}

		resp, err := client.Do(req)
		cancel()
		duration := time.Since(start)

		if err != nil {
			results <- CheckResult{
				TaskID: task.ID,
				URL:    task.URL,
				OK:     false,
				Err:    err,
			}
			continue
		}

		resp.Body.Close()

		ok := resp.StatusCode == http.StatusOK

		results <- CheckResult{
			TaskID:     task.ID,
			URL:        task.URL,
			OK:         ok,
			StatusCode: resp.StatusCode,
			Duration:   duration,
		}
	}
}

func main() {
	urls := []string{
		"https://example.org/",
		"https://example.org/",
		"https://example.org/",
		"https://example.org/",
		"https://example.org/",
		"https://example.org/",
		"https://example.org/",
		"https://example.org/",
	}

	jobs := make(chan CheckTask)
	results := make(chan CheckResult)

	runtime.GOMAXPROCS(1)

	var wg sync.WaitGroup

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for i, u := range urls {
			jobs <- CheckTask{
				ID:  i + 1,
				URL: u,
			}
		}
		close(jobs)
	}()

	for res := range results {
		if res.Err != nil {
			fmt.Printf("[TASK %02d] %s -> ERROR: %v \n", res.TaskID, res.URL, res.Err)
			continue
		}

		statusInfo := ""
		if res.OK {
			statusInfo = "OK (status 200)"
		} else {
			statusInfo = fmt.Sprintf("NOT OK (status %d)", res.StatusCode)
		}

		fmt.Printf("[TASK %02d] %s -> %s, time=%v \n", res.TaskID, res.URL, statusInfo, res.Duration)
	}
}
