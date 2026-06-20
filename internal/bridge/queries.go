package bridge

import (
	"database/sql"

	"inventory/internal/models"
	"inventory/internal/services"
)

// GetAlerts returns expiry/warranty alerts based on the configured warning threshold.
func (a *App) GetAlerts() ([]models.Alert, error) {
	return list(a, func(db *sql.DB) ([]models.Alert, error) {
		return services.GetAlerts(db, a.cfg.ExpiryWarningDays)
	})
}

// ListUsersForDropdown returns id and name pairs for all users.
func (a *App) ListUsersForDropdown() ([]models.DropdownItem, error) {
	return list(a, services.ListUsersForDropdown)
}

// ListDevicesForDropdown returns id, name and kind for all devices.
func (a *App) ListDevicesForDropdown() ([]models.DeviceDropdownItem, error) {
	return list(a, services.ListDevicesForDropdown)
}
