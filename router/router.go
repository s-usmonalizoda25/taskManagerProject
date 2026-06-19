package router

import (
	"net/http"

	"github.com/s-usmonalizoda25/taskManagerProject/handlers"
)

func NewRouter(taskHandler *handlers.TaskHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /tasks", taskHandler.GetTasks)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTaskByID)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	mux.HandleFunc("PUT /tasks/{id}/deactivate", taskHandler.DeactivateTask)

	return mux
}
