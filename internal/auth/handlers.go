package auth
import "bug-widget-saas/internal/middleware"


import (
    "crypto/rand"
    "database/sql"
    "encoding/hex"
    "encoding/json"
    "log"
    "net/http"
    "strings"
)

func generateSecureKey() (string, error) {
    bytes := make([]byte, 32) // 32 bytes = 64 hex chars
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return "sk_" + hex.EncodeToString(bytes), nil
}

func CreateTenant(db *sql.DB) http.HandlerFunc {
	
    return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        var req struct{ Name string }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        req.Name = strings.TrimSpace(req.Name)
        if len(req.Name) < 3 {
             http.Error(w, "Tenant name must be at least 3 characters", http.StatusBadRequest)
             return
        }

        apiKey, err := generateSecureKey()
        if err != nil {
             log.Printf("Error generating key: %v", err)
             http.Error(w, "Internal server error", http.StatusInternalServerError)
             return
        }
        
        _, err = db.Exec("INSERT INTO tenants (name, api_key) VALUES ($1, $2)", req.Name, apiKey)
        if err != nil {
            // Log the actual error for debugging
            log.Printf("Database error creating tenant: %v", err)
            // Return generic error to user to avoid leakage
            http.Error(w, "Could not create tenant. Name might already be taken.", http.StatusConflict)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(Tenant{Name: req.Name, APIKey: apiKey})
    })
}

func GetTenantByAPIKey(db *sql.DB) http.HandlerFunc {
    return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        if apiKey == "" {
            http.Error(w, "Missing X-API-Key", http.StatusUnauthorized)
            return
        }

        var tenant Tenant
        err := db.QueryRow("SELECT id, name, api_key FROM tenants WHERE api_key = $1", apiKey).
            Scan(&tenant.ID, &tenant.Name, &tenant.APIKey)
        if err == sql.ErrNoRows {
            // Use generic message to avoid enumerating valid keys vs invalid keys if distinguishable
            http.Error(w, "Invalid API key", http.StatusUnauthorized)
            return
        } else if err != nil {
            log.Printf("Database error checking key: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(tenant)
    })
}
