package main

import (
	"log"
	"net/http"
	"os"

	lofamhttp "github.com/stadtaev/lofam/backend/internal/http"
	"github.com/stadtaev/lofam/backend/internal/sqlite"
	"github.com/stadtaev/lofam/backend/internal/task"
)

func main() {
	dbPath := getEnv("DB_PATH", "lofam.db")
	port := getEnv("PORT", "8080")

	db, err := sqlite.New(dbPath)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	taskStore := sqlite.NewTaskStore(db)
	taskService := task.NewService(taskStore)
	server := lofamhttp.NewServer(taskService)

	log.Printf("starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, server.Router()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
