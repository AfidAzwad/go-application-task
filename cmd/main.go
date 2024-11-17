package main

import (
	"go-application-task/internal/routes"
	"log"
	"net/http"

	"go-application-task/pkg/db"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}

	defer func() {
		if err := db.WriteDB.Close(); err != nil {
			log.Printf("Error closing WriteDB: %v", err)
		}
	}()
	defer func() {
		if err := db.ReadDB.Close(); err != nil {
			log.Printf("Error closing ReadDB: %v", err)
		}
	}()

	router := routes.SetupRoutes()

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
