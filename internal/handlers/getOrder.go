package handlers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"go-application-task/internal/models"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Response struct to standardize the response format
type Response struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

// PaginatedResponse struct to handle paginated data
type PaginatedResponse struct {
	Data        interface{} `json:"data"`
	Total       int         `json:"total"`
	CurrentPage int         `json:"current_page"`
	PerPage     int         `json:"per_page"`
	TotalInPage int         `json:"total_in_page"`
	LastPage    int         `json:"last_page"`
}

// ListOrdersHandler handles the fetching of orders with pagination
func ListOrdersHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userID from token
		userID, err := GetUserIDFromToken(r, os.Getenv("JWT_SECRET"), db)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get pagination parameters from query params
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1 // default to page 1 if no valid page is provided
		}
		perPage, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || perPage < 1 {
			perPage = 10 // default to 10 per page
		}

		// Calculate offset for pagination
		offset := (page - 1) * perPage

		// Build SQL query with user ID filter and pagination
		query := `
			SELECT 
				store_id, 
				merchant_order_id, 
				recipient_name, 
				recipient_phone, 
				recipient_address, 
				recipient_city, 
				recipient_zone, 
				recipient_area, 
				delivery_type, 
				item_type, 
				transfer_status, 
				archive, 
				special_instruction, 
				item_quantity, 
				item_weight, 
				amount_to_collect, 
				item_description, 
				consignment_id, 
				order_status, 
				delivery_fee, 
				cod_fee, 
				user_id, 
				created_at
			FROM orders
			WHERE user_id=$1 AND transfer_status = 1 AND archive = 0
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`

		// Execute query with pagination
		rows, err := db.Queryx(query, userID, perPage, offset)
		if err != nil {
			log.Printf("Database query error: %v", err)
			http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Map results to an array of orders
		var orders []models.Order
		for rows.Next() {
			var order models.Order
			if err := rows.StructScan(&order); err != nil {
				log.Printf("Row scan error: %v", err)
				http.Error(w, "Failed to process orders", http.StatusInternalServerError)
				return
			}
			orders = append(orders, order)
		}

		// Count total orders to calculate pagination metadata
		var total int
		err = db.Get(&total, `
			SELECT COUNT(*) 
			FROM orders
			WHERE user_id=$1 AND transfer_status = 1 AND archive = 0
		`, userID)
		if err != nil {
			log.Printf("Error counting total orders: %v", err)
			http.Error(w, "Failed to calculate pagination", http.StatusInternalServerError)
			return
		}

		lastPage := (total / perPage)
		if total%perPage > 0 {
			lastPage++
		}

		paginatedResponse := PaginatedResponse{
			Data:        orders,
			Total:       total,
			CurrentPage: page,
			PerPage:     perPage,
			TotalInPage: len(orders),
			LastPage:    lastPage,
		}

		response := Response{
			Message: "Orders successfully fetched.",
			Type:    "success",
			Code:    200,
			Data:    paginatedResponse,
		}

		// Send the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode response: %v", err)
			http.Error(w, "Failed to process orders", http.StatusInternalServerError)
		}
	}
}
