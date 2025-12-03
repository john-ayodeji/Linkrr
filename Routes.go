package main

import (
	"net/http"

	analyticsHandler "github.com/john-ayodeji/Linkrr/internal/handlers/analyticsHandler"
	authHandler "github.com/john-ayodeji/Linkrr/internal/handlers/auth"
	"github.com/john-ayodeji/Linkrr/internal/handlers/redirectHandler"
	"github.com/john-ayodeji/Linkrr/internal/handlers/shortenerHandler"
	"github.com/john-ayodeji/Linkrr/internal/handlers/userHandler"
)

func RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/sign-up", authHandler.SignUp)
	mux.HandleFunc("POST /api/v1/auth/sign-up/", authHandler.SignUp)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/login/", authHandler.Login)
	mux.HandleFunc("GET /api/v1/auth/refresh", authHandler.RenewAccessToken)
	mux.HandleFunc("GET /api/v1/auth/refresh/", authHandler.RenewAccessToken)
	mux.HandleFunc("GET /api/v1/auth/revoke", authHandler.RevokeRefreshToken)
	mux.HandleFunc("GET /api/v1/auth/revoke/", authHandler.RevokeRefreshToken)
	mux.HandleFunc("POST /api/v1/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("POST /api/v1/auth/forgot-password/", authHandler.ForgotPassword)
	mux.HandleFunc("POST /api/v1/auth/reset-password", authHandler.ResetPassword)
	mux.HandleFunc("POST /api/v1/auth/reset-password/", authHandler.ResetPassword)
}

func RegisterShortenerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/shortener/new", shortenerHandler.HandleCreateUrl)
	mux.HandleFunc("POST /api/v1/shortener/new/", shortenerHandler.HandleCreateUrl)
	mux.HandleFunc("POST /api/v1/shortener/alias", shortenerHandler.CreateAlias)
	mux.HandleFunc("POST /api/v1/shortener/alias/", shortenerHandler.CreateAlias)
}

func RegisterRedirectRoute(mux *http.ServeMux) {
	mux.HandleFunc("GET /{urlCode}", redirectHandler.Redirect)
}

func RegisterUserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/links/me", userHandler.GetMyLinks)
	mux.HandleFunc("GET /api/v1/links/me/", userHandler.GetMyLinks)
}

func RegisterAnalyticsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/analytics/{urlCode}", analyticsHandler.GetURLAnalytics)
	mux.HandleFunc("GET /api/v1/analytics/{urlCode}/", analyticsHandler.GetURLAnalytics)
	mux.HandleFunc("GET /api/v1/analytics/{urlCode}/{alias}", analyticsHandler.GetAliasAnalytics)
	mux.HandleFunc("GET /api/v1/analytics/{urlCode}/{alias}/", analyticsHandler.GetAliasAnalytics)
	mux.HandleFunc("GET /api/v1/analytics/global", analyticsHandler.GetGlobalAnalytics)
	mux.HandleFunc("GET /api/v1/analytics/global/", analyticsHandler.GetGlobalAnalytics)
}
