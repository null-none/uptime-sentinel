package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
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

const workerCount = 10
const requestTimeout = 5 * time.Second

func worker(id int, jobs <-chan CheckTask, results chan<- CheckResult) {
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
				Err:    err,
			}
			continue
		}

		resp, err := client.Do(req)
		cancel()
		duration := time.Since(start)

		if err != nil {
			results <- CheckResult{
				TaskID:   task.ID,
				URL:      task.URL,
				OK:       false,
				Err:      err,
				Duration: duration,
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
		"https://google.com/",
	}

	jobs := make(chan CheckTask, 100)
	results := make(chan CheckResult, 100)

	runtime.GOMAXPROCS(1)

	for i := 1; i <= workerCount; i++ {
		go worker(i, jobs, results)
	}

	go func() {
		id := 0
		for {
			for _, u := range urls {
				id++
				jobs <- CheckTask{ID: id, URL: u}
			}
		}
	}()

	for {
		select {
		case res := <-results:
			if res.Err != nil {
				fmt.Printf("[TASK %06d] %s -> ERROR: %v, time=%v\n",
					res.TaskID, res.URL, res.Err, res.Duration)
				continue
			}

			status := "NOT OK"
			if res.OK {
				status = "OK"
			}

			fmt.Printf("[TASK %06d] %s -> %s (status %d, time=%v)\n",
				res.TaskID, res.URL, status, res.StatusCode, res.Duration)

		case <-time.After(5 * time.Second):
			fmt.Println("No results yet, service running...")
		}
	}
}
