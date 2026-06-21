package router

import (
	"net/http"

	"github.com/s-usmonalizoda25/taskManagerProject/handlers"
)

func NewRouter(taskHandler *handlers.TaskHandler, userHandler *handlers.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /tasks", taskHandler.GetTasks)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTaskByID)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)
	mux.HandleFunc("PATCH /tasks/{id}/status", taskHandler.PatchTaskStatus)
	mux.HandleFunc("PUT /tasks/{id}/deactivate", taskHandler.DeactivateTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /users", userHandler.GetUsers)
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)
	mux.HandleFunc("PUT /users/{id}/deactivate", userHandler.DeactivateUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)

	return mux
}
