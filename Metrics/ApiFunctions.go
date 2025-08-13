package Metrics

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func send_metrics(success bool, duration time.Duration) {
	// pseudo
	// if fail, add to failed requests
	// else, add time and request to window, to be calculated for average latency (ms)
	mu.Lock()
	defer mu.Unlock()

	if !success {
		failed_requets++
		return
	}

	time_window += duration
	requests_window++
	total_requests++
}

func ApiResponse(url string) ([]byte, error) {
	start_time := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error requesting api:", err)
		go send_metrics(false, 0)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response body:", err)
		go send_metrics(false, 0)
		return nil, err
	}

	duration := time.Since(start_time)

	// go routine to send data to cloud server
	go send_metrics(true, duration)

	return body, nil

}
