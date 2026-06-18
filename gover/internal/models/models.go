package models

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

type WindowsKey struct {
	ID         int     `json:"id"`
	LicenseKey string  `json:"licenseKey"`
	ComputerID *int    `json:"computerId"`
	Status     string  `json:"status"`
	Notes      *string `json:"notes"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

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

type User struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Surname   *string `json:"surname"`
	Status    string  `json:"status"`
	Notes     *string `json:"notes"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type AppConfig struct {
	DBPath            string `json:"dbPath"`
	Density           string `json:"density"`
	DarkMode          bool   `json:"darkMode"`
	ExpiryWarningDays int    `json:"expiryWarningDays"`
}

type Alert struct {
	Category      string `json:"category"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ExpiryDate    string `json:"expiryDate"`
	DaysRemaining int    `json:"daysRemaining"`
	Severity      string `json:"severity"`
}

type ComputerInput struct {
	Name           string  `json:"name"`
	Model          string  `json:"model"`
	UserID         *int    `json:"userId"`
	Status         string  `json:"status"`
	PurchaseDate   *string `json:"purchaseDate"`
	WarrantyExpiry *string `json:"warrantyExpiry"`
	Notes          *string `json:"notes"`
}

type SmartphoneInput struct {
	Name           string  `json:"name"`
	Model          string  `json:"model"`
	UserID         *int    `json:"userId"`
	Status         string  `json:"status"`
	PurchaseDate   *string `json:"purchaseDate"`
	WarrantyExpiry *string `json:"warrantyExpiry"`
	Notes          *string `json:"notes"`
}

type TabletInput struct {
	Name           string  `json:"name"`
	Model          string  `json:"model"`
	UserID         *int    `json:"userId"`
	Status         string  `json:"status"`
	PurchaseDate   *string `json:"purchaseDate"`
	WarrantyExpiry *string `json:"warrantyExpiry"`
	Notes          *string `json:"notes"`
}

type WindowsKeyInput struct {
	LicenseKey string  `json:"licenseKey"`
	ComputerID *int    `json:"computerId"`
	Status     string  `json:"status"`
	Notes      *string `json:"notes"`
}

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

type UserInput struct {
	Name    string  `json:"name"`
	Surname *string `json:"surname"`
	Status  string  `json:"status"`
	Notes   *string `json:"notes"`
}

type DropdownItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DeviceDropdownItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"` // "computer" | "smartphone" | "tablet"
}
