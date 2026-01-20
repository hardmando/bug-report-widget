package main

import (
	auth "bug-report-widget/internal/auth"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func connectDB() *sql.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	for i := 0; i < 30; i++ {
		db, err := sql.Open("postgres", connStr)
		if err == nil {
			if err := db.Ping(); err == nil {
				log.Println("Database connected")
				return db
			}
		}
		log.Printf("Database wait %d/30s: %v", i*2, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("Database timeout")
	return nil
}

func main() {
	db := connectDB()
	defer db.Close()

	http.HandleFunc("/register", auth.Register(db))
	http.HandleFunc("/login", auth.Login(db))
	http.HandleFunc("/github/login", auth.GitHubLogin(db))
	http.HandleFunc("/github/callback", auth.GitHubCallback(db))

	// Protected routes
	http.HandleFunc("/api-keys", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			auth.CreateAPIKey(db)(w, r)
		} else if r.Method == http.MethodGet {
			auth.GetAPIKeys(db)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// This is used by ingestion service or widget to validate key
	// But typically validation happens inside the service (internal validation)
	// or via an internal endpoint.
	// The previously existing /validate-key was for public? No, it's for internal use.
	// Let's keep /validate-key for `GetTenantByAPIKey`
	http.HandleFunc("/validate-key", auth.GetTenantByAPIKey(db))

	log.Println("Auth service starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
