package router

import (
	"net/http"
	"task-manager/internal/handler"
	"task-manager/internal/middleware"

	"github.com/gorilla/mux"
)

func SetupRouter(taskHandler *handler.TaskHandler) *mux.Router {
	router := mux.NewRouter()

	router.Use(addDefaultHeaders)

	public := router.PathPrefix("/public").Subrouter()
	public.HandleFunc("/tasks", taskHandler.GetTasks).Methods("GET")

	protected := router.PathPrefix("/tasks").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/", taskHandler.CreateTask).Methods("POST")

	return router
}

// Code has not been tested.
func addDefaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			rw.Header().Set("Access-Control-Allow-Origin", origin)
		}
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}
