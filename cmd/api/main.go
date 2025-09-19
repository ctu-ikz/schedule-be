package main

import (
	"context"
	"github.com/ctu-ikz/schedule-be/internal/api"
	"github.com/ctu-ikz/schedule-be/internal/api/handler"
	"github.com/ctu-ikz/schedule-be/internal/repository"
	"github.com/ctu-ikz/schedule-be/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("GOOSE_DBSTRING")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer pool.Close()

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	tokenRepo := repository.NewTokenRepository(pool)

	// Services
	authService := service.NewAuthService(userRepo, tokenRepo)

	// Handlers
	authHandler := handler.NewUserHandler(authService)

	// Router
	router := api.Router(authHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("🚀 API server running on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
