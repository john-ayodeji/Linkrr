package main

import (
	"net/http"
)

func (a *apiConfig) RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/sign-up", a.SignUp)
	mux.HandleFunc("POST /api/v1/auth/login", a.Login)
}