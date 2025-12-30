package http

import (
	"encoding/json"
	"net/http"

	"github.com/stadtaev/lofam/backend/internal/note"
)

func (s *Server) listNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := s.noteService.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, notes)
}

func (s *Server) createNote(w http.ResponseWriter, r *http.Request) {
	var req note.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	n, err := s.noteService.Create(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, n)
}

func (s *Server) getNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	n, err := s.noteService.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, n)
}

func (s *Server) updateNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	var req note.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	n, err := s.noteService.Update(r.Context(), id, req)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, n)
}

func (s *Server) deleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := s.noteService.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
