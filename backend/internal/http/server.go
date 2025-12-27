package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/stadtaev/lofam/backend/internal/task"
)

type Server struct {
	taskService *task.Service
}

func NewServer(taskService *task.Service) *Server {
	return &Server{taskService: taskService}
}

func (s *Server) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
	  AllowedOrigins:   []string{"http://localhost:3000"},
	  AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	  AllowedHeaders:   []string{"Accept", "Content-Type"},
	  AllowCredentials: false,
	  MaxAge:           300,
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Route("/api/tasks", func(r chi.Router) {
		r.Get("/", s.listTasks)
		r.Post("/", s.createTask)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", s.getTask)
			r.Put("/", s.updateTask)
			r.Delete("/", s.deleteTask)
		})
	})

	return r
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func handleError(w http.ResponseWriter, err error) {
	var validationErr task.ValidationError
	if errors.As(err, &validationErr) {
		writeError(w, http.StatusBadRequest, validationErr.Message)
		return
	}

	var notFoundErr task.NotFoundError
	if errors.As(err, &notFoundErr) {
		writeError(w, http.StatusNotFound, notFoundErr.Error())
		return
	}

	log.Printf("internal error: %v", err)
	writeError(w, http.StatusInternalServerError, "internal server error")
}

func parseID(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, task.ErrValidation("invalid id")
	}
	return id, nil
}
