package models

// Smartphone represents a mobile phone in the inventory.
type Smartphone struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Model          string  `json:"model"`
	UserID         *int    `json:"userId"`
	Status         string  `json:"status"`
	PurchaseDate   *string `json:"purchaseDate"`
	WarrantyExpiry *string `json:"warrantyExpiry"`
	Notes          *string `json:"notes"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

// SmartphoneInput carries the user-supplied fields for creating or updating a smartphone.
type SmartphoneInput struct {
	Name           string  `json:"name"`
	Model          string  `json:"model"`
	UserID         *int    `json:"userId"`
	Status         string  `json:"status"`
	PurchaseDate   *string `json:"purchaseDate"`
	WarrantyExpiry *string `json:"warrantyExpiry"`
	Notes          *string `json:"notes"`
}
