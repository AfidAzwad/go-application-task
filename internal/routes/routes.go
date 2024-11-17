package routes

import (
	"github.com/gorilla/mux"
	"go-application-task/internal/handlers"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// login route
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/refresh", handlers.RefreshTokenHandler).Methods("POST")

	return router
}
