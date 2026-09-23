package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"lab1/internal/model"
	"lab1/internal/service"
)

type PersonHandler struct {
	service *service.PersonService
}

func NewPersonHandler(service *service.PersonService) *PersonHandler {
	return &PersonHandler{service: service}
}

func NewRouter(h *PersonHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/persons", h.GetAll)
	mux.HandleFunc("POST /api/v1/persons", h.Create)
	mux.HandleFunc("GET /api/v1/persons/{id}", h.GetByID)
	mux.HandleFunc("PATCH /api/v1/persons/{id}", h.Update)
	mux.HandleFunc("DELETE /api/v1/persons/{id}", h.Delete)
	return mux
}

func (h *PersonHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeValidationError(w, "name", "invalid request body")
		return
	}

	id, err := h.service.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			writeValidationError(w, "name", "name is required")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Location", "/api/v1/persons/"+strconv.FormatInt(id, 10))
	w.WriteHeader(http.StatusCreated)
}

func (h *PersonHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, "person not found")
		return
	}

	person, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "person not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, person)
}

func (h *PersonHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	persons, err := h.service.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, persons)
}

func (h *PersonHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, "person not found")
		return
	}

	var req model.UpdatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeValidationError(w, "name", "invalid request body")
		return
	}

	person, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "person not found")
			return
		}
		if errors.Is(err, service.ErrValidation) {
			writeValidationError(w, "name", "name must not be empty")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, person)
}

func (h *PersonHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, "person not found")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "person not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, model.ErrorResponse{Message: message})
}

func writeValidationError(w http.ResponseWriter, field, message string) {
	writeJSON(w, http.StatusBadRequest, model.ValidationErrorResponse{
		Message: message,
		Errors:  map[string]string{field: message},
	})
}
