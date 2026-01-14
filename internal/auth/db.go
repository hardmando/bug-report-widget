package main

import (
    "database/sql"
    "log"
    "net/http"

    "bug-widget-saas/internal/ingestion"
)

func main() {
    db, _ := sql.Open("postgres", "postgres://postgres:devpass@postgres:5432/bugdb?sslmode=disable")
    defer db.Close()

    http.HandleFunc("/ingest/bugs", ingestion.IngestBug(db))
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Ingestion service running"))
    })

    log.Println("Ingestion service starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
