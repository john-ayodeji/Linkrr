package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/john-ayodeji/Linkrr/internal/config"
	"github.com/john-ayodeji/Linkrr/internal/database"
	"github.com/john-ayodeji/Linkrr/internal/services/auth"
	"github.com/john-ayodeji/Linkrr/internal/services/redirect"
	"github.com/john-ayodeji/Linkrr/internal/services/shortener"
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

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("DB open failed: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	cfg := &config.ApiConfig{
		Port:      3000,
		JWTSecret: jwtSecret,
		Db:        database.New(db),
	}

	authService.Cfg = cfg
	shortener.Cfg = cfg
	redirect.Cfg = cfg

	RegisterAuthRoutes(mux)
	RegisterShortenerRoutes(mux)
	RegisterRedirectRoute(mux)

	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	server := http.Server{Addr: addr, Handler: mux}
	log.Printf("Server started on port %d", cfg.Port)
	if err := http.ListenAndServe(server.Addr, server.Handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
