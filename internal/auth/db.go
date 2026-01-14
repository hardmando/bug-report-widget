package auth

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
	"os"
)

type Tenant struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    APIKey string `json:"api_key,omitempty"`
}

func ConnectDB() *sql.DB {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASS")
    dbname := os.Getenv("DB_NAME")

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    return db
}
