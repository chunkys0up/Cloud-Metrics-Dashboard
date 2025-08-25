package Metrics

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
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
	current_time_ms := strconv.FormatInt(time.Now().UnixMilli()-60000, 10)
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
