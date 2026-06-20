package bridge

import (
	"database/sql"

	"inventory/internal/models"
	"inventory/internal/services"
)

// ListUsers returns all users from the open database.
func (a *App) ListUsers() ([]models.User, error) { return list(a, services.ListUsers) }

// AddUser inserts a new user record.
func (a *App) AddUser(data models.UserInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertUser(db, data); return err })
}

// UpdateUser updates the user identified by id.
func (a *App) UpdateUser(id int, data models.UserInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateUser(db, id, data) })
}

// DeleteUser removes the user identified by id.
func (a *App) DeleteUser(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteUser(db, id) })
}
