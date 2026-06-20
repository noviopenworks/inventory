package bridge

import (
	"database/sql"

	"inventory/internal/models"
	"inventory/internal/services"
)

// ListComputers returns all computers from the open database.
func (a *App) ListComputers() ([]models.Computer, error) { return list(a, services.ListComputers) }

// AddComputer inserts a new computer record.
func (a *App) AddComputer(data models.ComputerInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertComputer(db, data); return err })
}

// UpdateComputer updates the computer identified by id.
func (a *App) UpdateComputer(id int, data models.ComputerInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateComputer(db, id, data) })
}

// DeleteComputer removes the computer identified by id.
func (a *App) DeleteComputer(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteComputer(db, id) })
}
