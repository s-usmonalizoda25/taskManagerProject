package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/s-usmonalizoda25/taskManagerProject/internal/models"
	"github.com/s-usmonalizoda25/taskManagerProject/internal/repository"
	"github.com/s-usmonalizoda25/taskManagerProject/internal/service"
)

type TaskHandler struct {
	service service.ITaskService
}

func NewTaskHandler(s service.ITaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var t models.Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.CreateTask(r.Context(), &t)
	if err != nil {
		if errors.Is(err, service.ErrEmptyTitle) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {

	status := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	tasks, err := h.service.GetTask(r.Context(), status, page, limit)
	if err != nil {

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTaskByID(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	var input models.Task
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateTask(r.Context(), uint(id), &input); err != nil {
		if errors.Is(err, service.ErrEmptyTitle) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, repository.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "task updated successfully"}`))
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTask(r.Context(), uint(id)); err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "task deleted successfully"}`))
}
