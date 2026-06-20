package bridge

import (
	"database/sql"

	"inventory/internal/models"
	"inventory/internal/services"
)

// ListSmartphones returns all smartphones from the open database.
func (a *App) ListSmartphones() ([]models.Smartphone, error) {
	return list(a, services.ListSmartphones)
}

// AddSmartphone inserts a new smartphone record.
func (a *App) AddSmartphone(data models.SmartphoneInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertSmartphone(db, data); return err })
}

// UpdateSmartphone updates the smartphone identified by id.
func (a *App) UpdateSmartphone(id int, data models.SmartphoneInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateSmartphone(db, id, data) })
}

// DeleteSmartphone removes the smartphone identified by id.
func (a *App) DeleteSmartphone(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteSmartphone(db, id) })
}
