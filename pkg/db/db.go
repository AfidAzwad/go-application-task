package db

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"go-application-task/configs"
)

var WriteDB *sqlx.DB
var ReadDB *sqlx.DB

// InitDB initializes db for read and write
func InitDB() error {
	// Write DB connection
	writeConfig := configs.GetWriteDBConfig()
	writeDSN := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		writeConfig.Host, writeConfig.Port, writeConfig.User, writeConfig.Password, writeConfig.DBName,
	)
	var err error
	WriteDB, err = sqlx.Connect("postgres", writeDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to write database: %w", err)
	}

	// Read DB connection
	readConfig := configs.GetReadDBConfig()
	readDSN := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		readConfig.Host, readConfig.Port, readConfig.User, readConfig.Password, readConfig.DBName,
	)
	ReadDB, err = sqlx.Connect("postgres", readDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to read database: %w", err)
	}

	log.Println("Databases connected successfully")

	// Apply migrations
	if err := ApplyMigrations(WriteDB); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// Call the seed function to ensure the default user is created
	seedDefaultUser()
	return nil
}

// seedDefaultUser ensures the default user is present in the database
func seedDefaultUser() {
	email := "01901901901@mailinator.com"
	password := "321dsa"

	// checking if the user already exists
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)"
	err := WriteDB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to check if default user exists: %v", err)
	}

	if !exists {
		// Hash
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash default user password: %v", err)
		}

		// Insert the default user
		insertQuery := "INSERT INTO users (email, password) VALUES ($1, $2)"
		_, err = WriteDB.Exec(insertQuery, email, string(hashedPassword))
		if err != nil {
			log.Fatalf("Failed to create default user: %v", err)
		}

		log.Println("Default user created successfully")
	} else {
		log.Println("Default user already exists")
	}
}
