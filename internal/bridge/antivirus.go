package bridge

import (
	"database/sql"

	"inventory/internal/models"
	"inventory/internal/services"
)

// ListAntivirus returns all antivirus records from the open database.
func (a *App) ListAntivirus() ([]models.Antivirus, error) { return list(a, services.ListAntivirus) }

// AddAntivirus inserts a new antivirus record.
func (a *App) AddAntivirus(data models.AntivirusInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertAntivirus(db, data); return err })
}

// UpdateAntivirus updates the antivirus record identified by id.
func (a *App) UpdateAntivirus(id int, data models.AntivirusInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateAntivirus(db, id, data) })
}

// DeleteAntivirus removes the antivirus record identified by id.
func (a *App) DeleteAntivirus(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteAntivirus(db, id) })
}
