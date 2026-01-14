package main

import (
    "log"
    "net/http"

    "bug-report-widget/internal/auth"
)

func main() {
    db := auth.ConnectDB()
    defer db.Close()

    http.HandleFunc("/tenants", auth.CreateTenant(db))
    http.HandleFunc("/validate-key", auth.GetTenantByAPIKey(db))
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Auth service running"))
    })

    log.Println("Auth service starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

