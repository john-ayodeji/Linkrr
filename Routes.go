package main

import (
	"net/http"

	"github.com/john-ayodeji/Linkrr/internal/handlers/auth"
)

func RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/sign-up", authHandler.SignUp)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/v1/auth/refresh", authHandler.RenewAccessToken)
	mux.HandleFunc("GET /api/v1/auth/revoke", authHandler.RevokeRefreshToken)
	mux.HandleFunc("POST /api/v1/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("POST /api/v1/auth/reset-password", authHandler.ResetPassword)
}
