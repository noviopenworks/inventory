// Package bridge exposes the Go backend to the Wails frontend via bound methods.
package bridge

import (
	"context"
	"database/sql"
	"errors"

	"inventory/internal/config"
	"inventory/internal/database"
	"inventory/internal/models"
)

// ErrNoDB is returned when no database is open.
var ErrNoDB = errors.New("no database open")

// list runs a read query, guarding the closed-db case and normalizing a nil
// slice to an empty one so the frontend always receives [] rather than null.
func list[T any](a *App, fn func(*sql.DB) ([]T, error)) ([]T, error) {
	if a.db == nil {
		return nil, ErrNoDB
	}
	out, err := fn(a.db)
	if out == nil {
		out = []T{}
	}
	return out, err
}

// withDB guards the closed-db case for write and side-effecting methods.
func (a *App) withDB(fn func(*sql.DB) error) error {
	if a.db == nil {
		return ErrNoDB
	}
	return fn(a.db)
}

// App holds the application state shared between the Go backend and the Wails
// frontend.
type App struct {
	db  *sql.DB
	cfg models.AppConfig
	ctx context.Context
}

// NewApp creates and returns a new App instance with default configuration.
func NewApp() *App {
	return &App{cfg: models.AppConfig{Density: "comfortable", ExpiryWarningDays: 30}}
}

// NewAppWithDB creates an App with a pre-opened database (used in tests).
func NewAppWithDB(db *sql.DB) *App {
	return &App{db: db, cfg: models.AppConfig{Density: "comfortable", ExpiryWarningDays: 30}}
}

// Startup is called by the Wails runtime when the application starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	cfg, err := config.Load()
	if err == nil {
		a.cfg = cfg
	} else {
		a.cfg = models.AppConfig{Density: "comfortable", DarkMode: false, ExpiryWarningDays: 30}
	}
	if a.cfg.DBPath != "" {
		if db, err := database.Open(a.cfg.DBPath); err == nil {
			if err := database.InitSchema(db); err == nil {
				a.db = db
			}
		}
	}
}
