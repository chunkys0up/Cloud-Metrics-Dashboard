package Metrics

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectToTracker(rdb *redis.Client, ctx context.Context) {
	MetricsCollected.RedisDB = rdb
	MetricsCollected.Ctx = ctx
}

func send_metrics(success bool, duration time.Duration) {
	// pseudo
	// if fail, add to failed requests
	// else, add time and request to window, to be calculated for average latency (ms)

	if !success {
		MetricsCollected.RedisDB.Incr(MetricsCollected.Ctx, "failed_requets")
		return
	}

	MetricsCollected.RedisDB.Incr(MetricsCollected.Ctx, "total_requests")

	// now store the duration with the timestamp as the ID
	err := MetricsCollected.RedisDB.XAdd(MetricsCollected.Ctx, &redis.XAddArgs{
		Stream: "time_window",
		ID:     "*",
		Values: map[string]any{
			"duration_ms": duration.Milliseconds(),
		},
	}).Err()

	if err != nil {
		panic(err)
	}
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
