package http

import (
	"encoding/json"
	"net/http"

	"github.com/stadtaev/lofam/backend/internal/wishlist"
)

func (s *Server) listWishlists(w http.ResponseWriter, r *http.Request) {
	items, err := s.wishlistService.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (s *Server) createWishlist(w http.ResponseWriter, r *http.Request) {
	var req wishlist.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := s.wishlistService.Create(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) getWishlist(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	item, err := s.wishlistService.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (s *Server) updateWishlist(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	var req wishlist.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := s.wishlistService.Update(r.Context(), id, req)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (s *Server) deleteWishlist(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := s.wishlistService.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
