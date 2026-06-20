package bridge

import (
	"database/sql"

	"inventory/internal/models"
	"inventory/internal/services"
)

// ListOtherSoftware returns all other-software records from the open database.
func (a *App) ListOtherSoftware() ([]models.OtherSoftware, error) {
	return list(a, services.ListOtherSoftware)
}

// AddOtherSoftware inserts a new other-software record.
func (a *App) AddOtherSoftware(data models.OtherSoftwareInput) error {
	return a.withDB(func(db *sql.DB) error {
		_, err := services.InsertOtherSoftware(db, data)
		return err
	})
}

// UpdateOtherSoftware updates the other-software record identified by id.
func (a *App) UpdateOtherSoftware(id int, data models.OtherSoftwareInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateOtherSoftware(db, id, data) })
}

// DeleteOtherSoftware removes the other-software record identified by id.
func (a *App) DeleteOtherSoftware(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteOtherSoftware(db, id) })
}
