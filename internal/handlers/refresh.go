package handlers

import (
	"encoding/json"
	"net/http"
	"os"

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

	accessToken, refreshToken, err := utils.GenerateToken(claims.Email, jwtSecret)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	// Return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
