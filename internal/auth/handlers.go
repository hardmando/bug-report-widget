package auth

import (
	"bug-report-widget/internal/middleware"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func init() {
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("default-secret-do-not-use-in-prod")
	}
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	TenantID string `json:"tenant_id"`
}

type Tenant struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	APIKey string `json:"api_key"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func generateSecureKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "sk_" + hex.EncodeToString(bytes), nil
}

// GenerateJWT creates a new token for a user
func GenerateJWT(userID, tenantID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"tenant_id": tenantID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func Register(db *sql.DB) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 1. Create Tenant (using email as name for simplicity, or handle better)
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		var tenantID string
		// Create a basic API key for the tenant initially
		apiKey, _ := generateSecureKey()
		err = tx.QueryRow("INSERT INTO tenants (name, api_key) VALUES ($1, $2) RETURNING id", req.Email, apiKey).Scan(&tenantID)
		if err != nil {
			log.Printf("Error creating tenant: %v", err)
			http.Error(w, "Email already registered", http.StatusConflict) // Assuming email is unique constraint on tenants(name) which might be true from previous schema
			return
		}

		// 2. Hash Password
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// 3. Create User
		var userID string
		err = tx.QueryRow("INSERT INTO users (email, password_hash, tenant_id) VALUES ($1, $2, $3) RETURNING id", req.Email, string(hashedPwd), tenantID).Scan(&userID)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "Transaction commit failed", http.StatusInternalServerError)
			return
		}

		// 4. Generate JWT
		token, _ := GenerateJWT(userID, tenantID)
		json.NewEncoder(w).Encode(TokenResponse{Token: token})
	})
}

func Login(db *sql.DB) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var userID, tenantID, hash string
		err := db.QueryRow("SELECT id, tenant_id, password_hash FROM users WHERE email = $1", req.Email).Scan(&userID, &tenantID, &hash)
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, _ := GenerateJWT(userID, tenantID)
		json.NewEncoder(w).Encode(TokenResponse{Token: token})
	})
}

// GitHub Auth Helpers
var (
	githubClientID     = os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
)

func GitHubLogin(db *sql.DB) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=user:email", githubClientID)
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	})
}

func GitHubCallback(db *sql.DB) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		// Exchange code for token
		// ... (Simplified for brevity, would normally make HTTP POST to GitHub)
		// For this mock implementation without real creds, we might fail or simulate.
		// But let's assume valid flow logic structure.

		// TODO: distinct implementation for real OAUTH call
		// For now, if we don't have secrets, we can't really do it.
		if githubClientID == "" {
			http.Error(w, "GitHub auth not configured", http.StatusInternalServerError)
			return
		}

		// Real implementation would go here:
		// 1. POST https://github.com/login/oauth/access_token
		// 2. GET https://api.github.com/user
		// 3. Find/Create user by github_id
		// 4. GenerateJWT

		http.Error(w, "Not fully implemented (requires valid GITHUB_CLIENT_ID env)", http.StatusNotImplemented)
	})
}

// Middleware to validate JWT
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing bearer token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Pass info maybe via context, skipping for simplicity in this swift implementation
		next(w, r)
	}
}

func CreateAPIKey(db *sql.DB) http.HandlerFunc {
	return middleware.CORS(AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// In a real app, we extract tenant_id from context/JWT claims
		// For now, let's parse the token again or rely on middleware setting ctx

		// Re-parsing for claims (inefficient but safe for now)
		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) { return jwtSecret, nil })
		claims := token.Claims.(jwt.MapClaims)
		tenantID := claims["tenant_id"].(string)

		apiKey, _ := generateSecureKey()
		// We might want to store multiple keys per tenant, but schema says api_key is on tenants table (single key?).
		// Oh, the first migration had `api_key` on `tenants` table.
		// If we want multiple keys, we need a separate `api_keys` table.
		// For now, let's just regenerate/update the key on the tenant.

		_, err := db.Exec("UPDATE tenants SET api_key = $1 WHERE id = $2", apiKey, tenantID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"api_key": apiKey})
	}))
}

func GetAPIKeys(db *sql.DB) http.HandlerFunc {
	return middleware.CORS(AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) { return jwtSecret, nil })
		claims := token.Claims.(jwt.MapClaims)
		tenantID := claims["tenant_id"].(string)

		var apiKey string
		err := db.QueryRow("SELECT api_key FROM tenants WHERE id = $1", tenantID).Scan(&apiKey)
		if err != nil {
			http.Error(w, "Error fetching key", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"api_key": apiKey})
	}))
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
