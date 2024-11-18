package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"go-application-task/internal/models"
	"go-application-task/pkg/db"
	"go-application-task/pkg/utils"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Constants for hardcoded values
const (
	ValidStoreID       = 131172
	ValidRecipientCity = 1
	ValidRecipientZone = 1
	ValidDeliveryType  = 48
	ValidItemType      = 2
	ValidItemQuantity  = 1
	ValidItemWeight    = 0.5
)

// GetUserIDFromToken extracts the user ID from the JWT token using the email
func GetUserIDFromToken(r *http.Request, jwtSecret string, db *sqlx.DB) (int, error) {
	// Get token from Authorization header (strip Bearer prefix)
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("authorization token is missing")
	}

	// Strip Bearer prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return 0, fmt.Errorf("authorization token format is incorrect")
	}

	// Parse the token and extract claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Return the secret key for validating the token
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("unable to parse claims")
	}

	// Extract email from claims
	email, ok := claims["email"].(string)
	if !ok {
		return 0, fmt.Errorf("email not found in token")
	}

	// Query the database to get the user_id for the given email
	var user models.User
	query := "SELECT * FROM users WHERE email = $1"
	err = db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("no user found with email %s", email)
	}
	if err != nil {
		return 0, fmt.Errorf("unable to retrieve user for email %s: %v", email, err)
	}
	return user.ID, nil
}

// ValidateOrderFields validates if the provided fields match the expected hardcoded values and if required fields are missing.
func ValidateOrderFields(order *models.Order) map[string][]string {
	errors := make(map[string][]string)

	// Validate store_id
	if order.StoreID == 0 {
		errors["store_id"] = append(errors["store_id"], "The store field is required")
	} else if order.StoreID != ValidStoreID {
		errors["store_id"] = append(errors["store_id"], "Wrong Store selected")
	}

	// Validate recipient_name
	if order.RecipientName == "" {
		errors["recipient_name"] = append(errors["recipient_name"], "The recipient name field is required.")
	}

	// Validate recipient_phone
	if order.RecipientPhone == "" {
		errors["recipient_phone"] = append(errors["recipient_phone"], "The recipient phone field is required.")
	}

	// Validate recipient_address
	if order.RecipientAddress == "" {
		errors["recipient_address"] = append(errors["recipient_address"], "The recipient address field is required.")
	}

	// Validate recipient_city
	if order.RecipientCity == 0 {
		errors["recipient_city"] = append(errors["recipient_city"], "The recipient city field is required.")
	} else if order.RecipientCity != ValidRecipientCity {
		errors["recipient_city"] = append(errors["recipient_city"], "Invalid city selected")
	}

	// Validate recipient_zone
	if order.RecipientZone == 0 {
		errors["recipient_zone"] = append(errors["recipient_zone"], "The recipient zone field is required.")
	} else if order.RecipientZone != ValidRecipientZone {
		errors["recipient_zone"] = append(errors["recipient_zone"], "Invalid zone selected")
	}

	// Validate delivery_type
	if order.DeliveryType == 0 {
		errors["delivery_type"] = append(errors["delivery_type"], "The delivery type field is required.")
	} else if order.DeliveryType != ValidDeliveryType {
		errors["delivery_type"] = append(errors["delivery_type"], "Invalid delivery type selected")
	}

	// Validate item_type
	if order.ItemType == 0 {
		errors["item_type"] = append(errors["item_type"], "The item type field is required.")
	} else if order.ItemType != ValidItemType {
		errors["item_type"] = append(errors["item_type"], "Invalid item type selected")
	}

	// Validate item_quantity
	if order.ItemQuantity == 0 {
		errors["item_quantity"] = append(errors["item_quantity"], "The item quantity field is required.")
	} else if order.ItemQuantity != ValidItemQuantity {
		errors["item_quantity"] = append(errors["item_quantity"], "Invalid item quantity selected")
	}

	// Validate item_weight
	if order.ItemWeight == 0 {
		errors["item_weight"] = append(errors["item_weight"], "The item weight field is required.")
	} else if order.ItemWeight != ValidItemWeight {
		errors["item_weight"] = append(errors["item_weight"], "Invalid item weight selected")
	}

	// Validate amount_to_collect
	if order.AmountToCollect == 0 {
		errors["amount_to_collect"] = append(errors["amount_to_collect"], "The amount to collect field is required.")
	}
	return errors
}

// validatePhone validates if the recipient phone number matches the Bangladesh phone number format
func validatePhone(phone string) bool {
	// Regex for Bangladesh phone number format (starting with 01 followed by 3-9, and then 8 digits)
	var phoneRegex = `^(01)[3-9]{1}[0-9]{8}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}

// CreateOrderHandler handles the creation of a new order
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	// Run user ID fetch in a goroutine
	userIDCh := make(chan int)
	go func() {
		userID, err := GetUserIDFromToken(r, jwtSecret, db.ReadDB)
		if err != nil {
			userIDCh <- -1 // signal error with a special value
		} else {
			userIDCh <- userID
		}
	}()

	// Validate required and hardcoded fields
	errors := ValidateOrderFields(&order)
	if len(errors) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Please fix the given errors",
			"type":    "error",
			"code":    422,
			"errors":  errors,
		})
		return
	}

	// Validate the recipient phone number
	if !validatePhone(order.RecipientPhone) {
		http.Error(w, "Invalid phone number", http.StatusBadRequest)
		return
	}

	// Retrieve user ID from goroutine
	userID := <-userIDCh
	if userID == -1 {
		http.Error(w, "Authentication error", http.StatusUnauthorized)
		return
	}

	// Generate consignment_id
	consignmentID, err := utils.GenerateConsignmentID("DA")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate consignment ID: %v", err), http.StatusInternalServerError)
		return
	}

	order.ConsignmentID = consignmentID
	order.OrderStatus = "pending"
	order.UserID = userID

	// Calculate delivery fee
	deliveryFee := 60
	if order.RecipientCity != ValidRecipientCity {
		deliveryFee = 100
	}
	order.DeliveryFee = float64(deliveryFee)
	codFee := 0.01 * order.AmountToCollect
	order.CODFee = codFee

	// Insert order into the database
	_, err = db.WriteDB.Exec(`
		INSERT INTO orders (store_id, recipient_name, recipient_phone, recipient_address, recipient_city, recipient_zone, recipient_area, delivery_type, item_type, item_quantity, item_weight, amount_to_collect, order_status, consignment_id, delivery_fee, cod_fee, user_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`, order.StoreID, order.RecipientName, order.RecipientPhone, order.RecipientAddress, order.RecipientCity, order.RecipientZone, order.RecipientArea, order.DeliveryType, order.ItemType, order.ItemQuantity, order.ItemWeight, order.AmountToCollect, order.OrderStatus, order.ConsignmentID, order.DeliveryFee, order.CODFee, order.UserID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create order: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Order Created Successfully",
		"type":    "success",
		"code":    200,
		"data": map[string]interface{}{
			"consignment_id":    consignmentID,
			"merchant_order_id": order.MerchantOrderID,
			"order_status":      order.OrderStatus,
			"delivery_fee":      deliveryFee,
		},
	})
}
