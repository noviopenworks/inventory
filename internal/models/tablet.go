package models

// Tablet represents a tablet device in the inventory.
type Tablet struct {
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

// TabletInput carries the user-supplied fields for creating or updating a tablet.
type TabletInput struct {
	Name           string  `json:"name"`
	Model          string  `json:"model"`
	UserID         *int    `json:"userId"`
	Status         string  `json:"status"`
	PurchaseDate   *string `json:"purchaseDate"`
	WarrantyExpiry *string `json:"warrantyExpiry"`
	Notes          *string `json:"notes"`
}
