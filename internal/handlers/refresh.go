package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"go-application-task/pkg/utils"
)

// RefreshTokenHandler handles the process of refreshing JWT tokens
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// inline struct
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// validating the refresh token and get the claims
	claims, err := utils.ValidateToken(req.RefreshToken, jwtSecret)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}
	accessTokenExpiry := time.Minute * 15 // Access token expires in 15 minutes
	refreshTokenExpiry := time.Hour * 24  // Refresh token expires in 1 day

	accessToken, refreshToken, err := utils.GenerateToken(claims.Email, jwtSecret, accessTokenExpiry, refreshTokenExpiry)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	// Return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":       accessToken,
		"refresh_token":      refreshToken,
		"expires_in":         strconv.Itoa(int(accessTokenExpiry.Seconds())),  // Expires in seconds for access token
		"refresh_expires_in": strconv.Itoa(int(refreshTokenExpiry.Seconds())), // Expires in seconds for refresh token
	})
}
