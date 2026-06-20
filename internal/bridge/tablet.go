package bridge

import (
	"database/sql"

	"inventory/internal/models"
	"inventory/internal/services"
)

// ListTablets returns all tablets from the open database.
func (a *App) ListTablets() ([]models.Tablet, error) { return list(a, services.ListTablets) }

// AddTablet inserts a new tablet record.
func (a *App) AddTablet(data models.TabletInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertTablet(db, data); return err })
}

// UpdateTablet updates the tablet identified by id.
func (a *App) UpdateTablet(id int, data models.TabletInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateTablet(db, id, data) })
}

// DeleteTablet removes the tablet identified by id.
func (a *App) DeleteTablet(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteTablet(db, id) })
}
