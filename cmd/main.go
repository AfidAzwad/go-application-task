package main

import (
	"github.com/jmoiron/sqlx"
	"go-application-task/internal/middleware"
	"go-application-task/internal/routes"
	"log"
	"net/http"

	"go-application-task/pkg/db"
)

// Helper function to close a database connection with proper error handling
func closeDatabaseConnection(dbConnection *sqlx.DB, dbName string) {
	if err := dbConnection.Close(); err != nil {
		log.Printf("Error closing %s: %v", dbName, err)
	} else {
		log.Printf("%s closed successfully", dbName)
	}
}

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}

	// Ensure proper closing of the database connections after the application exits
	defer closeDatabaseConnection(db.WriteDB, "WriteDB")
	defer closeDatabaseConnection(db.ReadDB, "ReadDB")

	router := routes.SetupRoutes()
	routerWithCors := middleware.EnableCors(router)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", routerWithCors))
}
