package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Kantha2004/SimpleJWT/internal/api"
	"github.com/Kantha2004/SimpleJWT/internal/auth"
	"github.com/Kantha2004/SimpleJWT/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	config := config.LoadConfig()

	jwtService := auth.NewJWTService(config.JWTSecret)

	routes := api.SetupRouters(jwtService)

	port := ":" + config.Port
	fmt.Printf("SimpleJWT server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, routes))
}
