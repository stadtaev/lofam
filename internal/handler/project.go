package handler

import (
	"encoding/json"
	"net/http"

	"lofam/internal/domain"
)

func (h *Handler) listProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectService.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, projects)
}

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	project, err := h.projectService.Create(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, project)
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	project, err := h.projectService.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, project)
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	var req domain.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	project, err := h.projectService.Update(r.Context(), id, req)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, project)
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := h.projectService.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listProjectTasks(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	tasks, err := h.taskService.ListByProjectID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}
