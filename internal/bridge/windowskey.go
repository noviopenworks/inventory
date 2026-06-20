package bridge

import (
	"database/sql"

	"inventory/internal/models"
	"inventory/internal/services"
)

// ListWindowsKeys returns all Windows license keys from the open database.
func (a *App) ListWindowsKeys() ([]models.WindowsKey, error) {
	return list(a, services.ListWindowsKeys)
}

// AddWindowsKey inserts a new Windows license key record.
func (a *App) AddWindowsKey(data models.WindowsKeyInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertWindowsKey(db, data); return err })
}

// UpdateWindowsKey updates the Windows license key identified by id.
func (a *App) UpdateWindowsKey(id int, data models.WindowsKeyInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateWindowsKey(db, id, data) })
}

// DeleteWindowsKey removes the Windows license key identified by id.
func (a *App) DeleteWindowsKey(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteWindowsKey(db, id) })
}
