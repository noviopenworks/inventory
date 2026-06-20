package models

// Antivirus represents an antivirus software license in the inventory.
type Antivirus struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	LicenseKey   string  `json:"licenseKey"`
	ComputerID   *int    `json:"computerId"`
	SmartphoneID *int    `json:"smartphoneId"`
	TabletID     *int    `json:"tabletId"`
	Status       string  `json:"status"`
	ExpiryDate   *string `json:"expiryDate"`
	Notes        *string `json:"notes"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

// AntivirusInput carries the user-supplied fields for creating or updating an antivirus record.
type AntivirusInput struct {
	Name         string  `json:"name"`
	LicenseKey   string  `json:"licenseKey"`
	ComputerID   *int    `json:"computerId"`
	SmartphoneID *int    `json:"smartphoneId"`
	TabletID     *int    `json:"tabletId"`
	Status       string  `json:"status"`
	ExpiryDate   *string `json:"expiryDate"`
	Notes        *string `json:"notes"`
}
