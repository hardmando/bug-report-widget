package ingestion

import (
	"bug-report-widget/internal/middleware"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type BugReport struct {
	Description string                   `json:"description"`
	URL         string                   `json:"url"`
	UserAgent   string                   `json:"userAgent"`
	Viewport    map[string]int           `json:"viewport"`
	Timestamp   string                   `json:"timestamp"`
	ConsoleLogs []map[string]interface{} `json:"consoleLogs"`
}

func IngestBug(db *sql.DB) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var bug BugReport
		json.NewDecoder(r.Body).Decode(&bug)

		log.Printf("Received bug from %s: %s", bug.URL, bug.Description[:50])

		// TODO: Validate API key → enqueue to Kafka
		// For now: just log + 200 OK

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"status": "ingested"})
	})
}
