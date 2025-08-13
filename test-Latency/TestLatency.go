package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/chunkys0up/Cloud-Metrics-Dashboard/Metrics"
)

type SampleData struct {
	Likes    uint
	Dislikes uint
	Views    uint32
}

func sampleHTTPFunction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Stats := SampleData {
		Likes: 1034,
		Dislikes: 87,
		Views: 234781,
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
	mux.HandleFunc("/getData/", http.HandlerFunc(sampleHTTPFunction))

	// starting the server
	fmt.Println("Server starting at http://localhost:8080")
	go http.ListenAndServe(":8080", mux)

	// simulate 10 calls
	for i := 1; i <= 10;i++ {
		fmt.Println(i)
		Metrics.ApiResponse("http://localhost:8080/getData")
	}
}
