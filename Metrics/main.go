package Metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// global metrics
var (
	mu              sync.Mutex
	total_requests  int
	failed_requets  int
	cpu_usage       float64
	bytesRecvRate   float64
	bytesSentRate   float64
	average_latency float64
	time_window     time.Duration
	requests_window int
)

type SiteData struct {
	TotalRequests    int
	FailedRequests   int
	AverageLatencyMs float64
}

type NetworkTraffic struct {
	RxBytesRate float64
	TxBytesRate float64
}

type ServerData struct {
	CpuUsed        float64
	MemoryUsed     float64
	DiskUsed       float64
	NetworkTraffic NetworkTraffic
}

type Metrics struct {
	SiteData   SiteData
	ServerData ServerData
}

type Report struct {
	Timestamp string
	Metrics   Metrics
}

// --- API
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

	memory_used := sampleMemory()
	disk_used := sampleDisk()

	report := Report{
		Timestamp: time.Now().Format("2006-01-02T15:04:05"),
		Metrics: Metrics{
			SiteData: SiteData{
				TotalRequests:    total_requests,
				FailedRequests:   failed_requets,
				AverageLatencyMs: average_latency,
			},
			ServerData: ServerData{
				CpuUsed:    cpu_usage,
				MemoryUsed: memory_used,
				DiskUsed:   disk_used,
				NetworkTraffic: NetworkTraffic{
					RxBytesRate: bytesRecvRate,
					TxBytesRate: bytesSentRate,
				},
			},
		},
	}

	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	w.Write(jsonData)
}

func main() {
	go sampleCPU()
	go sampleBytes()
	go sampleLatency()

	mux := http.NewServeMux()
	mux.HandleFunc("/get/", http.HandlerFunc(ServeHTTP))

	// run the server
	fmt.Println("Server starting at http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
