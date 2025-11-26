package main

import (
	"net/http"
)

func (a *apiConfig) RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/sign-up", a.SignUp)
	mux.HandleFunc("POST /api/v1/auth/login", a.Login)
	mux.HandleFunc("GET /api/v1/auth/refresh", a.RenewAccessToken)
	mux.HandleFunc("GET /api/v1/auth/revoke", a.RevokeRefreshToken)
}
