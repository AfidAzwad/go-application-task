package models

// Credentials represents the user-provided login details.
type Credentials struct {
	Email    string `json:"email"`    // Maps to "email" in JSON
	Password string `json:"password"` // Maps to "password" in JSON
}
