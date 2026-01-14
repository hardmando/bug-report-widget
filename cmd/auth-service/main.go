package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/validate-key", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "X-API-Key")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        apiKey := r.Header.Get("X-API-Key")
        if apiKey == "dev_api_key_123" {
            w.Header().Set("Content-Type", "application/json")
            fmt.Fprint(w, `{"id":"dev","name":"dev-tenant","api_key":"dev_api_key_123"}`)
            return
        }
        http.Error(w, "Invalid API key", http.StatusUnauthorized)
    })
    
    http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
        fmt.Fprint(w, "Auth service running")
    })
    
    log.Println("Auth service starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
