// Package models defines the data types shared between the backend and frontend.
package models

// Computer represents a desktop or laptop computer in the inventory.
type Computer struct {
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

// User represents a person who can be assigned devices or software.
type User struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Surname   *string `json:"surname"`
	Status    string  `json:"status"`
	Notes     *string `json:"notes"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

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

// ComputerInput carries the user-supplied fields for creating or updating a computer.
type ComputerInput struct {
	Name           string  `json:"name"`
	Model          string  `json:"model"`
	UserID         *int    `json:"userId"`
	Status         string  `json:"status"`
	PurchaseDate   *string `json:"purchaseDate"`
	WarrantyExpiry *string `json:"warrantyExpiry"`
	Notes          *string `json:"notes"`
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

// WindowsKeyInput carries the user-supplied fields for creating or updating a Windows key.
type WindowsKeyInput struct {
	LicenseKey string  `json:"licenseKey"`
	ComputerID *int    `json:"computerId"`
	Status     string  `json:"status"`
	Notes      *string `json:"notes"`
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

// UserInput carries the user-supplied fields for creating or updating a user.
type UserInput struct {
	Name    string  `json:"name"`
	Surname *string `json:"surname"`
	Status  string  `json:"status"`
	Notes   *string `json:"notes"`
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
