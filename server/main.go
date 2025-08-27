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
	AverageLatency float32
	CpuUsed        float64
	MemoryUsed     float64
	DiskUsed       float64
	RxBytesRate    float32
	TxBytesRate    float32
}

// ---API CALL---
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Unblocks CORS to host: 5173
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	report := Report{
		Timestamp:      time.Now().Format("2006-01-02T15:04:05"),
		TotalRequests:  120,
		FailedRequests: 13,
		AverageLatency: 9.4,
		CpuUsed:        67.8,
		MemoryUsed:     45.8,
		DiskUsed:       12.3,
		RxBytesRate:    1200.4,
		TxBytesRate:    604.3,
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
	mux.HandleFunc("/get/", http.HandlerFunc(ServeHTTP))

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx = context.Background()

	// Share redis client and context
	// created ctx and redis client for adding stream
	// on the client side, will create another rdb so we can send data to db that will be added to stream
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
			Metrics.addStream()
		}
	}()

	<-quit
	// fmt.Println("\nStream entries")
	// entries, err := rdb.XRange(ctx, "mystream", "-", "+").Result()
	// if err != nil {
	// 	panic(err)
	// }

	// for index, entry := range entries {
	// 	fmt.Printf("%d) 1) \"%s\"\n", index, entry.ID)
	// 	i := 2
	// 	for k, v := range entry.Values {
	// 		fmt.Printf("   %d) \"%s\"\n", i, k)
	// 		fmt.Printf("   %d) \"%v\"\n", i+1, v)
	// 		i += 2
	// 	}
	// }

	// Clears database
	err := rdb.FlushDB(ctx).Err()
	if err != nil {
		panic(err)
	}
	fmt.Print("\nSuccessfully cleared database\n")
}
