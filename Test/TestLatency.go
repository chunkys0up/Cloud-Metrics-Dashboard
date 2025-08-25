package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chunkys0up/Cloud-Metrics-Dashboard/Metrics"
	"github.com/redis/go-redis/v9"
	"net/http"
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

	// 3) Calls at random times with random amounts at a time (fix)
	for i := 1; i <= 25; i++ {
		fmt.Println("Call:", i)
		Metrics.ApiResponse("http://localhost:8081/getData")
	}

	
}
