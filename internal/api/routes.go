package api

import (
	"net/http"

	"github.com/Kantha2004/SimpleJWT/internal/auth"
)

// Helper function for public routes with logging
func registerPublicRoute(mux *http.ServeMux, path string, handler http.HandlerFunc) {
	mux.Handle(path, LoggingMiddleware(http.HandlerFunc(handler)))
}

// Helper function for protected routes with logging + JWT
func registerProtectedRoute(mux *http.ServeMux, path string, handler http.HandlerFunc, jwtService *auth.JWTService) {
	mux.Handle(path, Chain(
		http.HandlerFunc(handler),
		LoggingMiddleware,
		JWTMiddleware(jwtService),
	))
}

func SetupRouters(jwtService *auth.JWTService) *http.ServeMux {
	mux := http.NewServeMux()

	// Public routes
	publicRoutes := map[string]http.HandlerFunc{
		"/ping": PingHandler,
	}

	for path, handler := range publicRoutes {
		registerPublicRoute(mux, path, handler)
	}

	return mux
}
