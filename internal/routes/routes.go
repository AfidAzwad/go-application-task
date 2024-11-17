package routes

import (
	"github.com/gorilla/mux"
	"go-application-task/internal/handlers"
	"go-application-task/internal/middleware"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/refresh", handlers.RefreshTokenHandler).Methods("POST")

	createOrderRoute := router.HandleFunc("/create_order", handlers.CreateOrderHandler).Methods("POST")
	createOrderRoute.Handler(middleware.JWTMiddleware(createOrderRoute.GetHandler()))

	return router
}
