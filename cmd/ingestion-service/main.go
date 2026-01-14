package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
)

type BugReport struct {
    Description string `json:"description"`
}

func main() {
    http.HandleFunc("/ingest/bugs", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "X-API-Key, Content-Type")
        
        if r.Method == "OPTIONS" { return }
        
        body, _ := io.ReadAll(r.Body)
        log.Printf("Bug ingested: %s", body)
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"status": "ingested"})
    })
    
    http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
        fmt.Fprint(w, "Ingestion service running")
    })
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
