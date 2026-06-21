package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/s-usmonalizoda25/taskManagerProject/internal/models"
	"github.com/s-usmonalizoda25/taskManagerProject/internal/service"
	"github.com/s-usmonalizoda25/taskManagerProject/pkg/logger"
)

type UserHandler struct {
	service service.IUserService
	log     *logger.Logger
}

func NewUserHandler(s service.IUserService, log *logger.Logger) *UserHandler {
	return &UserHandler{service: s, log: log}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateUser(r.Context(), &u); err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user ID format", http.StatusBadRequest)
		return
	}

	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateUser(r.Context(), uint(id), &input); err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "user updated successfully"}`))
}

func (h *UserHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user ID format", http.StatusBadRequest)
		return
	}

	if err := h.service.DeactivateUser(r.Context(), uint(id)); err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "user deactivated successfully (soft-deleted)"}`))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user ID format", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteUser(r.Context(), uint(id)); err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "user hard-deleted from database"}`))
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	allStr := r.URL.Query().Get("all")

	var page, limit int
	var err error

	if pageStr != "" {
		if page, err = strconv.Atoi(pageStr); err != nil || page <= 0 {
			http.Error(w, "invalid page parameter", http.StatusBadRequest)
			return
		}
	}
	if limitStr != "" {
		if limit, err = strconv.Atoi(limitStr); err != nil || limit <= 0 {
			http.Error(w, "invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	var users []models.User
	if allStr == "true" {
		users, err = h.service.GetAllUsers(r.Context(), page, limit)
	} else {
		users, err = h.service.GetActiveUsers(r.Context(), page, limit)
	}
	if err != nil {
		HandleError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
