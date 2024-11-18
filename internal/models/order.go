package models

import "time"

// Order represents the order data structure
type Order struct {
	StoreID            int       `json:"store_id" validate:"required" db:"store_id"`
	MerchantOrderID    *string   `json:"merchant_order_id,omitempty" db:"merchant_order_id"`
	RecipientName      string    `json:"recipient_name" validate:"required" db:"recipient_name"`
	RecipientPhone     string    `json:"recipient_phone" validate:"required,phone" db:"recipient_phone"`
	RecipientAddress   string    `json:"recipient_address" validate:"required" db:"recipient_address"`
	RecipientCity      int       `json:"recipient_city" validate:"required" db:"recipient_city"`
	RecipientZone      int       `json:"recipient_zone" validate:"required" db:"recipient_zone"`
	RecipientArea      int       `json:"recipient_area" validate:"required" db:"recipient_area"`
	DeliveryType       int       `json:"delivery_type" validate:"required" db:"delivery_type"`
	ItemType           int       `json:"item_type" validate:"required" db:"item_type"`
	TransferStatus     int       `json:"transfer_status,omitempty" db:"transfer_status"`
	Archive            int       `json:"archive,omitempty" db:"archive"`
	SpecialInstruction *string   `json:"special_instruction,omitempty" db:"special_instruction"`
	ItemQuantity       int       `json:"item_quantity" validate:"required" db:"item_quantity"`
	ItemWeight         float64   `json:"item_weight" validate:"required" db:"item_weight"`
	AmountToCollect    float64   `json:"amount_to_collect" validate:"required" db:"amount_to_collect"`
	ItemDescription    *string   `json:"item_description,omitempty" db:"item_description"`
	ConsignmentID      string    `json:"consignment_id" validate:"required,len=16" db:"consignment_id"`
	OrderStatus        string    `json:"order_status" validate:"required,oneof=pending completed cancelled" db:"order_status"`
	DeliveryFee        float64   `json:"delivery_fee" validate:"required" db:"delivery_fee"`
	CODFee             float64   `json:"cod_fee" validate:"required" db:"cod_fee"`
	UserID             int       `json:"user_id" db:"user_id"`
	OrderCreatedAt     time.Time `json:"created_at" db:"created_at"`
}
