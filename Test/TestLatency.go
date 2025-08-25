package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/chunkys0up/Cloud-Metrics-Dashboard/Metrics"
	"github.com/redis/go-redis/v9"
)

type SampleData struct {
	Likes    uint
	Dislikes uint
	Views    uint32
}

func sampleHTTPFunction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Stats := SampleData{
		Likes:    1034,
		Dislikes: 87,
		Views:    234781,
	}

	jsonData, err := json.MarshalIndent(Stats, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	w.Write(jsonData)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/getData", http.HandlerFunc(sampleHTTPFunction))

	// starting the server
	fmt.Println("Server starting at http://localhost:8081")
	go http.ListenAndServe(":8081", mux)

	// Initialization
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx := context.Background()
	Metrics.ConnectToTracker(rdb, ctx)

	// --- Tests ---
	// 1) Verify redis client exists
	fmt.Println("redis client:", Metrics.MetricsCollected.RedisDB)

	// 2) Test Api response through custom function that tracks the speed of call
	fmt.Print("Testing api response\n")
	resp, err := Metrics.ApiResponse("http://localhost:8081/getData")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp))

	// 3) Calls 4 batches with random ammounts to make sure latency works
	// does not count 404 errors as a failed request because it still "successfully" something
	for i := range 4 {
		fmt.Println("Batch:", i + 1)
		m := rand.Intn(20) + 5

		for j := range m {
			if j % 3 == 0 {
				Metrics.ApiResponse("http://localhost:8082/failRequest")
			} else {
				Metrics.ApiResponse("http://localhost:8081/getData")
			}
		}

		time.Sleep(3 * time.Second)
	}
}
