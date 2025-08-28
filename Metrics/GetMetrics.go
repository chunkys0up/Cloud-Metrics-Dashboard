package Metrics

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"strconv"
	"sync"
	"time"
)

var MetricsCollected struct {
	mu      sync.Mutex
	RedisDB *redis.Client
	Ctx     context.Context

	CpuUsed       float64
	BytesRecvRate float64
	BytesSentRate float64
}

func SampleMemory() float64 {
	v, _ := mem.VirtualMemory()

	return v.UsedPercent
}

func SampleDisk() float64 {
	v, _ := disk.Usage("/")

	return v.UsedPercent
}

// get the timed cpu usage every second
func SampleCPU() {
	for {
		percent, err := cpu.Percent(time.Second, false)
		if err == nil && len(percent) > 0 {
			MetricsCollected.mu.Lock()
			MetricsCollected.CpuUsed = percent[0]
			MetricsCollected.mu.Unlock()
		}
	}
}

// Times the Bytes sent and received and calculates the rate
func SampleBytes() {
	kilobytes_per_second := float64(1) / 1000

	for {
		initialIoData, err := net.IOCounters(false)
		if err != nil {
			fmt.Println("Error getting intial I/O data:", err)
			return
		}

		initialTime := time.Now()

		time.Sleep(1 * time.Second)

		endIoData, err := net.IOCounters(false)
		if err != nil {
			fmt.Println("Error getting end I/O data:", err)
			return
		}

		endTime := time.Now()

		duration := endTime.Sub(initialTime).Seconds()
		bytesRecv := float64(endIoData[0].BytesRecv - initialIoData[0].BytesRecv)
		bytesSent := float64(endIoData[0].BytesSent - initialIoData[0].BytesSent)

		MetricsCollected.mu.Lock()
		MetricsCollected.BytesRecvRate = bytesRecv / duration * float64(kilobytes_per_second)
		MetricsCollected.BytesSentRate = bytesSent / duration * float64(kilobytes_per_second)
		MetricsCollected.mu.Unlock()
	}
}

// Upates the window and updates the average latency
func SampleLatency() float64 {
	current_time_ms := strconv.FormatInt(time.Now().UnixMilli()-10000, 10)
	_, err := MetricsCollected.RedisDB.XTrimMinID(MetricsCollected.Ctx, "time_window", current_time_ms).Result()
	if err != nil {
		panic(err)
	}

	// get response amounts
	request_window, err := MetricsCollected.RedisDB.XLen(MetricsCollected.Ctx, "time_window").Result()
	if err != nil {
		panic(err)
	}

	// get the entries in the window
	entries, err := MetricsCollected.RedisDB.XRange(MetricsCollected.Ctx, "time_window", "-", "+").Result()
	if err != nil {
		panic(err)
	}

	// calculate duration sum
	var total_time int64
	for _, entry := range entries {
		raw := entry.Values["duration_ms"].(string)
		converted, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			fmt.Println("Failed to convert:", err)
			return 0
		}
		total_time += converted
	}

	if request_window > 0 {
		average_latency := float64(total_time) / float64(request_window)
		return average_latency
	}

	return 0
}

func getRedisValue(key string) int64 {
	raw, err := MetricsCollected.RedisDB.Get(MetricsCollected.Ctx, key).Result()
	if err != nil {
		return 0
	}

	return ToInt64(raw)
}

// create stream that adds metrics to redis
func AddStream() {
	total_requests := getRedisValue("total_requests")
	failed_requests := getRedisValue("failed_requests")

	latency_ms := SampleLatency()
	memory_used := SampleMemory()
	disk_used := SampleDisk()

	id, err := MetricsCollected.RedisDB.XAdd(MetricsCollected.Ctx, &redis.XAddArgs{
		Stream: "mystream",
		Values: map[string]any{
			"total_requests":     total_requests,
			"failed_requests":    failed_requests,
			"average_latency_ms": latency_ms,
			"cpu_usage":          MetricsCollected.CpuUsed,
			"memory_usage":       memory_used,
			"disk_usage":         disk_used,
			"rx_bps":             MetricsCollected.BytesRecvRate,
			"tx_bps":             MetricsCollected.BytesSentRate,
		},
	}).Result()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Latest id: %s\n", id)
}

// converts raw string to int64
func ToInt64(raw string) int64 {
	i64, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		fmt.Print("Failed to convert to Int 64\n")
		return 0
	}

	return i64
}

func ToFloat64(raw string) float64 {
	f64, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		fmt.Print("Failed to convert to Float 64\n")
		return 0
	}

	return f64
}