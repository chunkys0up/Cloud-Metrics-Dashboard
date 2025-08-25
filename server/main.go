package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/chunkys0up/Cloud-Metrics-Dashboard/Metrics"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx context.Context

func getRedisValue(key string) int64 {
	raw, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return 0
	}

	converted, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0
	}

	return converted
}

// create stream that adds metrics to redis
func addStream() {
	total_requests := getRedisValue("total_requests")
	failed_requests := getRedisValue("failed_requests")

	latency_ms := Metrics.SampleLatency()
	memory_used := Metrics.SampleMemory()
	disk_used := Metrics.SampleDisk()

	id, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "mystream",
		Values: map[string]any{
			"total_requests":     total_requests,
			"failed_requests":    failed_requests,
			"average_latency_ms": latency_ms,
			"cpu_usage":          Metrics.MetricsCollected.CpuUsed,
			"memory_usage":       memory_used,
			"disk_usage":         disk_used,
			"rx_bps":             Metrics.MetricsCollected.BytesRecvRate,
			"tx_bps":             Metrics.MetricsCollected.BytesSentRate,
		},
	}).Result()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Latest id: %s\n", id)
}

func main() {
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
			addStream()
		}
	}()

	// Clear database
	<-quit
	fmt.Println("\nStream entries")
	entries, err := rdb.XRange(ctx, "mystream", "-", "+").Result()
	if err != nil {
		panic(err)
	}

	for index, entry := range entries {
		fmt.Printf("%d) 1) \"%s\"\n", index, entry.ID)
		i := 2
		for k, v := range entry.Values {
			fmt.Printf("   %d) \"%s\"\n", i, k)
			fmt.Printf("   %d) \"%v\"\n", i+1, v)
			i += 2
		}
	}

	err = rdb.FlushDB(ctx).Err()
	if err != nil {
		panic(err)
	}
	fmt.Print("\nSuccessfully cleared database\n")
}
