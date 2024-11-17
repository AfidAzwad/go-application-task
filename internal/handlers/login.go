package handlers

import (
	"encoding/json"
	"go-application-task/internal/models"
	"go-application-task/pkg/db"
	"go-application-task/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

// LoginHandler handles user login requests. It authenticates the user based on their email and password,
// and generates JWT tokens upon successful authentication.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// decoding credentials
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// checking user
	var storedPassword string
	query := "SELECT password FROM users WHERE email=$1"
	err = db.ReadDB.QueryRow(query, creds.Email).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// generating a JWT token for the user
	accessToken, refreshToken, err := utils.GenerateToken(creds.Email, jwtSecret)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
