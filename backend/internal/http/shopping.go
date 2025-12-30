package http

import (
	"encoding/json"
	"net/http"

	"github.com/stadtaev/lofam/backend/internal/shopping"
)

func (s *Server) listShoppingItems(w http.ResponseWriter, r *http.Request) {
	items, err := s.shoppingService.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (s *Server) createShoppingItem(w http.ResponseWriter, r *http.Request) {
	var req shopping.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := s.shoppingService.Create(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) deleteShoppingItem(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := s.shoppingService.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
