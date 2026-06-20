package models

// WindowsKey represents a Windows operating system license key.
type WindowsKey struct {
	ID         int     `json:"id"`
	LicenseKey string  `json:"licenseKey"`
	ComputerID *int    `json:"computerId"`
	Status     string  `json:"status"`
	Notes      *string `json:"notes"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

// WindowsKeyInput carries the user-supplied fields for creating or updating a Windows key.
type WindowsKeyInput struct {
	LicenseKey string  `json:"licenseKey"`
	ComputerID *int    `json:"computerId"`
	Status     string  `json:"status"`
	Notes      *string `json:"notes"`
}
