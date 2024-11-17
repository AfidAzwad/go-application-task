package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// GenerateToken generates both access and refresh tokens
func GenerateToken(email, secret string, accessTokenExpiry time.Duration, refreshTokenExpiry time.Duration) (accessToken string, refreshToken string, err error) {
	// Access token (short expiry)
	accessClaims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenExpiry).Unix(),
		},
	}

	accessTokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenJwt.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	// Refresh token (long expiry)
	refreshClaims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenExpiry).Unix(),
		},
	}

	refreshTokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenJwt.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	// Parse the token with claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	// Check if there was an error parsing the token
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Verify if the token is valid and if the claims can be asserted
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Println("Token is valid, claims:", claims)
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
