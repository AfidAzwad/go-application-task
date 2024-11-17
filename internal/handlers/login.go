package handlers

import (
	"encoding/json"
	"go-application-task/internal/models"
	"go-application-task/pkg/db"
	"go-application-task/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
		// User not found, return error response
		response := map[string]interface{}{
			"message": "The user credentials were incorrect.",
			"type":    "error",
			"code":    400,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized) // Set the status code
		json.NewEncoder(w).Encode(response)    // Send the JSON response
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	accessTokenExpirySecondStr := os.Getenv("ACCESS_TOKEN_EXPIRY_SECOND")
	refreshTokenExpirySecondStr := os.Getenv("REFRESH_TOKEN_EXPIRY_SECOND")

	if accessTokenExpirySecondStr == "" {
		accessTokenExpirySecondStr = "3600" // Default to 1 hour
	}
	if refreshTokenExpirySecondStr == "" {
		refreshTokenExpirySecondStr = "7200" // Default to 2 hours
	}

	// Convert the environment variables to integers
	accessTokenExpirySecond, err := strconv.Atoi(accessTokenExpirySecondStr)
	if err != nil {
		log.Fatalf("Error parsing ACCESS_TOKEN_EXPIRY_SECOND: %v", err)
	}

	refreshTokenExpirySecond, err := strconv.Atoi(refreshTokenExpirySecondStr)
	if err != nil {
		log.Fatalf("Error parsing REFRESH_TOKEN_EXPIRY_SECOND: %v", err)
	}

	// Convert the seconds into time.Duration
	accessTokenExpiry := time.Second * time.Duration(accessTokenExpirySecond)
	refreshTokenExpiry := time.Second * time.Duration(refreshTokenExpirySecond)

	// generating a JWT token for the user
	accessToken, refreshToken, err := utils.GenerateToken(creds.Email, jwtSecret, accessTokenExpiry, refreshTokenExpiry)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":       accessToken,
		"refresh_token":      refreshToken,
		"access_expires_in":  strconv.Itoa(int(accessTokenExpiry.Seconds())),  // Expires in seconds for access token
		"refresh_expires_in": strconv.Itoa(int(refreshTokenExpiry.Seconds())), // Expires in seconds for refresh token
	})
}
