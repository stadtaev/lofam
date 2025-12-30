package main

import (
	"log"
	"net/http"
	"os"

	lofamhttp "github.com/stadtaev/lofam/backend/internal/http"
	"github.com/stadtaev/lofam/backend/internal/note"
	"github.com/stadtaev/lofam/backend/internal/sqlite"
	"github.com/stadtaev/lofam/backend/internal/task"
	"github.com/stadtaev/lofam/backend/internal/wishlist"
)

func main() {
	dbPath := getEnv("DB_PATH", "lofam.db")
	port := getEnv("PORT", "8080")
	staticDir := getEnv("STATIC_DIR", "./static")

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

	noteStore := sqlite.NewNoteStore(db)
	noteService := note.NewService(noteStore)

	wishlistStore := sqlite.NewWishlistStore(db)
	wishlistService := wishlist.NewService(wishlistStore)

	server := lofamhttp.NewServer(taskService, noteService, wishlistService, staticDir)

	log.Printf("starting server on :%s", port)
	log.Printf("serving static files from %s", staticDir)
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
