package config

import "github.com/john-ayodeji/Linkrr/internal/database"

type ApiConfig struct {
	Port      int
	Db        *database.Queries
	JWTSecret string
}
