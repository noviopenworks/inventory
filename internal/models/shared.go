// Package models defines the data types shared between the backend and frontend.
package models

// AppConfig holds the persisted application settings.
type AppConfig struct {
	DBPath            string `json:"dbPath"`
	Density           string `json:"density"`
	DarkMode          bool   `json:"darkMode"`
	ExpiryWarningDays int    `json:"expiryWarningDays"`
}

// Alert represents an expiry or warranty alert for an inventory item.
type Alert struct {
	Category      string `json:"category"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ExpiryDate    string `json:"expiryDate"`
	DaysRemaining int    `json:"daysRemaining"`
	Severity      string `json:"severity"`
}

// DropdownItem is a minimal id+name pair used in select dropdowns.
type DropdownItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DeviceDropdownItem extends DropdownItem with a device kind discriminator.
type DeviceDropdownItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"` // "computer" | "smartphone" | "tablet"
}
