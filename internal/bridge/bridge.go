// Package bridge exposes the Go backend to the Wails frontend via bound methods.
package bridge

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"inventory/internal/backup"
	"inventory/internal/config"
	"inventory/internal/database"
	"inventory/internal/models"
	"inventory/internal/services"
)

var errNoDB = errors.New("no database open")

// list runs a read query, guarding the closed-db case and normalizing a nil
// slice to an empty one so the frontend always receives [] rather than null.
func list[T any](a *App, fn func(*sql.DB) ([]T, error)) ([]T, error) {
	if a.db == nil {
		return nil, errNoDB
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
		return errNoDB
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

// ListComputers returns all computers from the open database.
func (a *App) ListComputers() ([]models.Computer, error) { return list(a, services.ListComputers) }

// ListSmartphones returns all smartphones from the open database.
func (a *App) ListSmartphones() ([]models.Smartphone, error) {
	return list(a, services.ListSmartphones)
}

// ListTablets returns all tablets from the open database.
func (a *App) ListTablets() ([]models.Tablet, error) { return list(a, services.ListTablets) }

// ListWindowsKeys returns all Windows license keys from the open database.
func (a *App) ListWindowsKeys() ([]models.WindowsKey, error) {
	return list(a, services.ListWindowsKeys)
}

// ListAntivirus returns all antivirus records from the open database.
func (a *App) ListAntivirus() ([]models.Antivirus, error) { return list(a, services.ListAntivirus) }

// ListOtherSoftware returns all other-software records from the open database.
func (a *App) ListOtherSoftware() ([]models.OtherSoftware, error) {
	return list(a, services.ListOtherSoftware)
}

// ListUsers returns all users from the open database.
func (a *App) ListUsers() ([]models.User, error) { return list(a, services.ListUsers) }

// GetAlerts returns expiry/warranty alerts based on the configured warning threshold.
func (a *App) GetAlerts() ([]models.Alert, error) {
	return list(a, func(db *sql.DB) ([]models.Alert, error) {
		return services.GetAlerts(db, a.cfg.ExpiryWarningDays)
	})
}

// GetDatabasePath returns the filesystem path of the currently open database.
func (a *App) GetDatabasePath() (string, error) {
	return a.cfg.DBPath, nil
}

// GetConfig returns the current application configuration.
func (a *App) GetConfig() (models.AppConfig, error) {
	return a.cfg, nil
}

// SetConfig persists the supplied configuration and updates the in-memory copy.
func (a *App) SetConfig(cfg models.AppConfig) error {
	a.cfg = cfg
	return config.Save(cfg)
}

// NewDatabase creates a new SQLite database at path and opens it.
func (a *App) NewDatabase(path string) error {
	db, err := database.Open(path)
	if err != nil {
		return err
	}
	if err := database.InitSchema(db); err != nil {
		_ = db.Close()
		return err
	}
	if a.db != nil {
		_ = a.db.Close()
	}
	a.db = db
	a.cfg.DBPath = path
	return config.Save(a.cfg)
}

// OpenDatabase opens the existing SQLite database at path.
func (a *App) OpenDatabase(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("file not found: %s", path)
	}
	return a.NewDatabase(path)
}

// NewDatabaseDialog opens a native save-file dialog and creates a new database
// at the chosen path. Returns nil if the user cancels. The dialog itself
// requires a live Wails context and is exercised by the smoke test; the
// underlying create logic lives in NewDatabase.
func (a *App) NewDatabaseDialog() error {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "inventory.db",
		Filters:         []runtime.FileFilter{{DisplayName: "SQLite database", Pattern: "*.db"}},
	})
	if err != nil {
		return err
	}
	if path == "" {
		return nil // user cancelled
	}
	return a.NewDatabase(path)
}

// OpenDatabaseDialog opens a native open-file dialog and opens the selected
// database. Returns nil if the user cancels. The dialog itself requires a live
// Wails context and is exercised by the smoke test; the underlying open logic
// lives in OpenDatabase.
func (a *App) OpenDatabaseDialog() error {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{{DisplayName: "SQLite database", Pattern: "*.db"}},
	})
	if err != nil {
		return err
	}
	if path == "" {
		return nil // user cancelled
	}
	return a.OpenDatabase(path)
}

// ExportCSV builds a CSV for the given category and prompts the user to save it.
func (a *App) ExportCSV(category string) error {
	return a.withDB(func(db *sql.DB) error {
		rows, err := services.BuildCSV(db, category)
		if err != nil {
			return err
		}
		path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
			DefaultFilename: category + ".csv",
			Filters:         []runtime.FileFilter{{DisplayName: "CSV", Pattern: "*.csv"}},
		})
		if err != nil {
			return err
		}
		if path == "" {
			return nil // user cancelled
		}
		return services.WriteCSV(path, rows)
	})
}

// ---------------------------------------------------------------------------
// Computer write methods
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Smartphone write methods
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Tablet write methods
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// WindowsKey write methods
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Antivirus write methods
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// OtherSoftware write methods
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// User write methods
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Dropdown helpers
// ---------------------------------------------------------------------------

// ListUsersForDropdown returns id and name pairs for all users.
func (a *App) ListUsersForDropdown() ([]models.DropdownItem, error) {
	return list(a, services.ListUsersForDropdown)
}

// ListDevicesForDropdown returns id, name and kind for all devices.
func (a *App) ListDevicesForDropdown() ([]models.DeviceDropdownItem, error) {
	return list(a, services.ListDevicesForDropdown)
}

var _ = backup.DB
