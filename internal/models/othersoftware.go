package models

// OtherSoftware represents a software license that is not antivirus.
type OtherSoftware struct {
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

// OtherSoftwareInput carries the user-supplied fields for creating or updating an other-software record.
type OtherSoftwareInput struct {
	Name         string  `json:"name"`
	LicenseKey   string  `json:"licenseKey"`
	ComputerID   *int    `json:"computerId"`
	SmartphoneID *int    `json:"smartphoneId"`
	TabletID     *int    `json:"tabletId"`
	Status       string  `json:"status"`
	ExpiryDate   *string `json:"expiryDate"`
	Notes        *string `json:"notes"`
}
