package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go-application-task/pkg/utils"
)

// RefreshTokenHandler handles the process of refreshing JWT tokens
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Inline struct for decoding request
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Decode the incoming request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Get the JWT secret from environment variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	claims, err := utils.ValidateToken(req.RefreshToken, jwtSecret)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
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

	accessToken, refreshToken, err := utils.GenerateToken(claims.Email, jwtSecret, accessTokenExpiry, refreshTokenExpiry)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	expiresInSeconds := int(accessTokenExpiry.Seconds())
	refreshExpiresInSeconds := int(refreshTokenExpiry.Seconds())

	// Return the generated tokens and expiry times
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":       accessToken,
		"refresh_token":      refreshToken,
		"expires_in":         strconv.Itoa(expiresInSeconds),
		"refresh_expires_in": strconv.Itoa(refreshExpiresInSeconds),
	})
}
