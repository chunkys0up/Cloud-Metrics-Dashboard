package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/chunkys0up/Cloud-Metrics-Dashboard/Metrics"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx context.Context

type Report struct {
	Timestamp      string
	TotalRequests  int64
	FailedRequests int64
	AverageLatency float64
	CpuUsed        float64
	MemoryUsed     float64
	DiskUsed       float64
	RxBytesRate    float64
	TxBytesRate    float64
}

// ---API CALL---
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Unblocks CORS to host: 5173
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	stream, err := rdb.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"mystream", "$"},
		Block: 0,
		Count: 1,
	}).Result()

	if err != nil {
		panic(err)
	}

	latest := stream[0].Messages[0].Values

	report := Report{
		Timestamp:      time.Now().Format("2006-01-02T15:04:05"),
		TotalRequests:  Metrics.ToInt64(latest["total_requests"].(string)),
		FailedRequests: Metrics.ToInt64(latest["failed_requests"].(string)),
		AverageLatency: Metrics.ToFloat64(latest["average_latency_ms"].(string)),
		CpuUsed:        Metrics.ToFloat64(latest["cpu_usage"].(string)),
		MemoryUsed:     Metrics.ToFloat64(latest["memory_usage"].(string)),
		DiskUsed:       Metrics.ToFloat64(latest["disk_usage"].(string)),
		RxBytesRate:    Metrics.ToFloat64(latest["rx_bps"].(string)),
		TxBytesRate:    Metrics.ToFloat64(latest["tx_bps"].(string)),
	}

	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	w.Write(jsonData)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/getData", http.HandlerFunc(ServeHTTP))

	fmt.Println("Server starting at http://localhost:8080")
    go http.ListenAndServe(":8080", mux)

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx = context.Background()

	// Share redis client and context
	// created ctx and redis client for adding stream
	// on the client side, will create another rdb so we can send data to the redis that will be added to stream
	Metrics.ConnectToTracker(rdb, ctx)

	fmt.Print("Metric Sampling starting...\n")
	go Metrics.SampleCPU()
	go Metrics.SampleBytes()

	// Channel to listen for OS signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// adds to redis every second
	go func() {
		for {
			time.Sleep(1 * time.Second)
			Metrics.AddStream()
		}
	}()

	<-quit

	// Clears database
	err := rdb.FlushDB(ctx).Err()
	if err != nil {
		panic(err)
	}
	fmt.Print("\nSuccessfully cleared database\n")
}
