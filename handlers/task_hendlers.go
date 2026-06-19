package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/s-usmonalizoda25/taskManagerProject/internal/models"
	"github.com/s-usmonalizoda25/taskManagerProject/internal/service"
	"github.com/s-usmonalizoda25/taskManagerProject/pkg/errs"
	"github.com/s-usmonalizoda25/taskManagerProject/pkg/logger"
	"go.uber.org/zap"
)

type TaskHandler struct {
	service service.ITaskService
	log     *logger.Logger
}

func NewTaskHandler(s service.ITaskService, log *logger.Logger) *TaskHandler {
	return &TaskHandler{service: s, log: log}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.log.Error("Failed to decode create task body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.service.CreateTask(r.Context(), &t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	allStr := r.URL.Query().Get("all")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	var tasks []models.Task
	var err error

	if allStr == "true" {
		tasks, err = h.service.GetAllTasks(r.Context(), status, page, limit)
	} else {
		tasks, err = h.service.GetActiveTasks(r.Context(), status, page, limit)
	}

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
		http.Error(w, "invalid task ID format", http.StatusBadRequest)
		return
	}
	task, err := h.service.GetTaskByID(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid task ID format", http.StatusBadRequest)
		return
	}
	var input models.Task
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.service.UpdateTask(r.Context(), uint(id), &input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "task updated successfully"}`))
}

func (h *TaskHandler) DeactivateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid task ID format", http.StatusBadRequest)
		return
	}
	if err := h.service.DeactivateTask(r.Context(), uint(id)); err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "task deactivated successfully (soft-deleted)"}`))
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid task ID format", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteTask(r.Context(), uint(id)); err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "task hard-deleted from database"}`))
}

