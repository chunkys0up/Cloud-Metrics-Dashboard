package main

import (
    "fmt"
    "net/http"
    "github.com/chunkys0up/Cloud-Metrics-Dashboard/Metrics"
)

func main() {
    go Metrics.SampleCPU()       // assuming exported function from Metrics package
    go Metrics.SampleBytes()
    go Metrics.SampleLatency()

    mux := http.NewServeMux()
    mux.HandleFunc("/get/", Metrics.ServeHTTP)

    fmt.Println("Server starting at http://localhost:8080")
    http.ListenAndServe(":8080", mux)
}
