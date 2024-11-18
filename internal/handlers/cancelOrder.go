package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go-application-task/internal/models"
	"log"
	"net/http"
	"os"
)

// CancelOrderHandler handles the cancellation of an order
func CancelOrderHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract consignment_id from query parameters
		consignmentID := r.URL.Query().Get("consignment_id")
		if consignmentID == "" {
			http.Error(w, "Consignment ID is required", http.StatusBadRequest)
			return
		}

		// Extract userID from token
		userID, err := GetUserIDFromToken(r, os.Getenv("JWT_SECRET"), db)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the order exists and retrieve it
		var order models.Order
		query := `SELECT order_status, user_id FROM orders WHERE consignment_id = $1 AND user_id = $2`
		err = db.Get(&order, query, consignmentID, userID)
		if err != nil {
			log.Printf("Order retrieval error: %v", err)
			http.Error(w, "Failed to fetch order", http.StatusInternalServerError)
			return
		}

		// Check if order status is not "Cancelled"
		if order.OrderStatus == "cancelled" {
			http.Error(w, "Order already cancelled", http.StatusConflict)
			return
		}

		// Check if order status is "Completed"
		if order.OrderStatus != "pending" {
			http.Error(w, "Please contact cx to cancel order", http.StatusConflict)
			return
		}

		// Update the order status to "Cancelled"
		updateQuery := `UPDATE orders SET order_status = 'cancelled' WHERE consignment_id = $1 AND user_id = $2`
		_, err = db.Exec(updateQuery, consignmentID, userID)
		if err != nil {
			log.Printf("Failed to update order status: %v", err)
			http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
			return
		}

		// Send success response
		response := Response{
			Message: fmt.Sprintf("Order with consignment ID %s successfully cancelled.", consignmentID),
			Type:    "success",
			Code:    200,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
