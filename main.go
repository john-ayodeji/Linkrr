package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/john-ayodeji/Linkrr/internal/config"
	"github.com/john-ayodeji/Linkrr/internal/database"
	"github.com/john-ayodeji/Linkrr/internal/events_workers"
	"github.com/john-ayodeji/Linkrr/internal/services/analytics"
	"github.com/john-ayodeji/Linkrr/internal/services/auth"
	"github.com/john-ayodeji/Linkrr/internal/services/redirect"
	"github.com/john-ayodeji/Linkrr/internal/services/shortener"
	"github.com/john-ayodeji/Linkrr/internal/services/users"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	port      int
	db        *database.Queries
	jwtSecret string
}

func main() {
	mux := http.NewServeMux()
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set in environment")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set in environment")
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 3030 // Default port
	}

	IpstackApiKey := os.Getenv("IPSTACK_API_KEY")
	IpStackURL := os.Getenv("IPSTACK_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("DB open failed: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	cfg := &config.ApiConfig{
		Port:          port,
		JWTSecret:     jwtSecret,
		Db:            database.New(db),
		IpStackApiKey: IpstackApiKey,
		IpStackURl:    IpStackURL,
	}

	authService.Cfg = cfg
	shortener.Cfg = cfg
	redirect.Cfg = cfg
	analytics.Cfg = cfg
	users.Cfg = cfg

	for i := 0; i < 5; i++ {
		go events_workers.SignUpEmailWorker(authService.SignUpEvent)
		go events_workers.LoginEmailWorker(authService.LoginEvent)
		go events_workers.ForgotPasswordEmailWorker(authService.ForgotPasswordEvent)
		go events_workers.ChangedPasswordEmailWorker(authService.ResetPasswordEvent)

		go analytics.GetClickData(redirect.RedirectEvent)
	}

	go analytics.AggregateAnalytics(analytics.AnalyticsEvent)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	RegisterAuthRoutes(mux)
	RegisterShortenerRoutes(mux)
	RegisterRedirectRoute(mux)
	RegisterUserRoutes(mux)
	RegisterAnalyticsRoutes(mux)

	// Add CORS and panic recovery middleware
	handler := corsMiddleware(panicRecoveryMiddleware(mux))

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	server := http.Server{Addr: addr, Handler: handler}
	log.Printf("Server started on port %d", cfg.Port)
	if err := http.ListenAndServe(server.Addr, server.Handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Panic recovery middleware
func panicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"Internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
