package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/john-ayodeji/Linkrr/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	port int
	db *database.Queries
	jwtSecret string
}

func main() {
	mux := http.NewServeMux()
	var cfg apiConfig
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

	cfg.db = database.New(db)
	cfg.jwtSecret = jwtSecret
	cfg.port = 3000

	cfg.RegisterAuthRoutes(mux)

	addr := fmt.Sprintf("localhost:%d", cfg.port)
	server := http.Server{Addr: addr, Handler: mux}
	log.Printf("Server started on port %d", cfg.port)
	if err := http.ListenAndServe(server.Addr, server.Handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}