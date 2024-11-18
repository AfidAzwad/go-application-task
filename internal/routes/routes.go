package routes

import (
	"github.com/gorilla/mux"
	"go-application-task/internal/handlers"
	"go-application-task/internal/middleware"
	"go-application-task/pkg/db"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/refresh", handlers.RefreshTokenHandler).Methods("POST")

	createOrderRoute := router.HandleFunc("/create_order", handlers.CreateOrderHandler).Methods("POST")
	createOrderRoute.Handler(middleware.JWTMiddleware(createOrderRoute.GetHandler()))

	getOrderRoute := router.HandleFunc("/orders", handlers.ListOrdersHandler(db.ReadDB)).Methods("GET")
	getOrderRoute.Handler(middleware.JWTMiddleware(handlers.ListOrdersHandler(db.ReadDB)))

	cancelOrderRoute := router.HandleFunc("/cancel-order", handlers.CancelOrderHandler(db.WriteDB)).Methods("POST")
	cancelOrderRoute.Handler(middleware.JWTMiddleware(handlers.CancelOrderHandler(db.WriteDB)))

	return router
}
