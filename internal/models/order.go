package models

// Order represents the order data structure
type Order struct {
	StoreID            int     `json:"store_id" validate:"required"`
	MerchantOrderID    *string `json:"merchant_order_id,omitempty"`
	RecipientName      string  `json:"recipient_name" validate:"required"`
	RecipientPhone     string  `json:"recipient_phone" validate:"required,phone"`
	RecipientAddress   string  `json:"recipient_address" validate:"required"`
	RecipientCity      int     `json:"recipient_city" validate:"required"`
	RecipientZone      int     `json:"recipient_zone" validate:"required"`
	RecipientArea      int     `json:"recipient_area" validate:"required"`
	DeliveryType       int     `json:"delivery_type" validate:"required"`
	ItemType           int     `json:"item_type" validate:"required"`
	SpecialInstruction *string `json:"special_instruction,omitempty"`
	ItemQuantity       int     `json:"item_quantity" validate:"required"`
	ItemWeight         float64 `json:"item_weight" validate:"required"`
	AmountToCollect    float64 `json:"amount_to_collect" validate:"required"`
	ItemDescription    *string `json:"item_description,omitempty"`
	ConsignmentID      string  `json:"consignment_id" validate:"required,len=16"`
	OrderStatus        string  `json:"order_status" validate:"required,oneof=pending completed cancelled"`
	DeliveryFee        float64 `json:"delivery_fee" validate:"required"`
	CODFee             float64 `json:"cod_fee" validate:"required"`
	UserID             int     `json:"user_id"`
}
