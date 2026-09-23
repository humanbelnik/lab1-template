package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"lab1/internal/config"
	"lab1/internal/handler"
	"lab1/internal/repository"
	"lab1/internal/service"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := repository.Migrate(ctx, pool); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	repo := repository.NewPostgresPersonRepository(pool)
	svc := service.NewPersonService(repo)
	h := handler.NewPersonHandler(svc)
	router := handler.NewRouter(h)

	log.Printf("server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
