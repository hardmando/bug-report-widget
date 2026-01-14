package ingestion

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
    "os"
)

func ConnectDB() *sql.DB {
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), 
        os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    return db
}
