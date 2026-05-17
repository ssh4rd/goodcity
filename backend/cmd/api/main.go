package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"goodcity/backend/internal/config"
	"goodcity/backend/internal/db"
	"goodcity/backend/internal/handler"
	"goodcity/backend/internal/repository/postgres"
	"goodcity/backend/internal/service"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	// Run migrations
	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	// Connect pool
	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	// Repositories
	userRepo := postgres.NewUserRepository(pool)
	practiceRepo := postgres.NewPracticeRepository(pool)
	commentRepo := postgres.NewCommentRepository(pool)
	ratingRepo := postgres.NewRatingRepository(pool)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	practiceSvc := service.NewPracticeService(practiceRepo, commentRepo)
	ratingSvc := service.NewRatingService(ratingRepo, practiceRepo)

	// Router
	router := handler.NewRouter(authSvc, practiceSvc, ratingSvc)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func runMigrations(databaseURL string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open db for migrations: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	// MIGRATIONS_DIR can be set to override; defaults to ../migrations (relative to backend/ cwd)
	migrationsDir := "../migrations"
	if v := os.Getenv("MIGRATIONS_DIR"); v != "" {
		migrationsDir = v
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
