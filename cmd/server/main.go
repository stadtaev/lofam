package main

import (
	"log"
	"net/http"
	"os"

	"lofam/internal/handler"
	"lofam/internal/repository/sqlite"
	"lofam/internal/service"
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

	projectRepo := sqlite.NewProjectRepository(db)
	taskRepo := sqlite.NewTaskRepository(db)

	projectService := service.NewProjectService(projectRepo)
	taskService := service.NewTaskService(taskRepo, projectRepo)

	h := handler.New(projectService, taskService)

	log.Printf("starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, h.Router()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
