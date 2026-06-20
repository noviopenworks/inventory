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

type App struct {
	db  *sql.DB
	cfg models.AppConfig
	ctx context.Context
}

func NewApp() *App {
	return &App{cfg: models.AppConfig{Density: "comfortable", ExpiryWarningDays: 30}}
}

// NewAppWithDB creates an App with a pre-opened database (used in tests).
func NewAppWithDB(db *sql.DB) *App {
	return &App{db: db, cfg: models.AppConfig{Density: "comfortable", ExpiryWarningDays: 30}}
}

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

func (a *App) ListComputers() ([]models.Computer, error) {
	if a.db == nil {
		return nil, errNoDB
	}
	out, err := services.ListComputers(a.db)
	if out == nil {
		out = []models.Computer{}
	}
	return out, err
}

func (a *App) ListSmartphones() ([]models.Smartphone, error) {
	if a.db == nil {
		return nil, errNoDB
	}
	out, err := services.ListSmartphones(a.db)
	if out == nil {
		out = []models.Smartphone{}
	}
	return out, err
}

func (a *App) ListTablets() ([]models.Tablet, error) {
	if a.db == nil {
		return nil, errNoDB
	}
	out, err := services.ListTablets(a.db)
	if out == nil {
		out = []models.Tablet{}
	}
	return out, err
}

func (a *App) ListWindowsKeys() ([]models.WindowsKey, error) {
	if a.db == nil {
		return nil, errNoDB
	}
	out, err := services.ListWindowsKeys(a.db)
	if out == nil {
		out = []models.WindowsKey{}
	}
	return out, err
}

func (a *App) ListAntivirus() ([]models.Antivirus, error) {
	if a.db == nil {
		return nil, errNoDB
	}
	out, err := services.ListAntivirus(a.db)
	if out == nil {
		out = []models.Antivirus{}
	}
	return out, err
}

func (a *App) ListOtherSoftware() ([]models.OtherSoftware, error) {
	if a.db == nil {
		return nil, errNoDB
	}
	out, err := services.ListOtherSoftware(a.db)
	if out == nil {
		out = []models.OtherSoftware{}
	}
	return out, err
}

func (a *App) ListUsers() ([]models.User, error) {
	if a.db == nil {
		return nil, errNoDB
	}
	out, err := services.ListUsers(a.db)
	if out == nil {
		out = []models.User{}
	}
	return out, err
}

func (a *App) GetAlerts() ([]models.Alert, error) {
	if a.db == nil {
		return nil, errNoDB
	}
	return services.GetAlerts(a.db, a.cfg.ExpiryWarningDays)
}

func (a *App) GetDatabasePath() (string, error) {
	return a.cfg.DBPath, nil
}

func (a *App) GetConfig() (models.AppConfig, error) {
	return a.cfg, nil
}

func (a *App) SetConfig(cfg models.AppConfig) error {
	a.cfg = cfg
	return config.Save(cfg)
}

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

func (a *App) ExportCSV(category string) error {
	if a.db == nil {
		return errNoDB
	}
	rows, err := services.BuildCSV(a.db, category)
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
}

// ---------------------------------------------------------------------------
// Computer write methods
// ---------------------------------------------------------------------------

func (a *App) AddComputer(data models.ComputerInput) error {
	if a.db == nil {
		return errNoDB
	}
	_, err := services.InsertComputer(a.db, data)
	return err
}

func (a *App) UpdateComputer(id int, data models.ComputerInput) error {
	if a.db == nil {
		return errNoDB
	}
	return services.UpdateComputer(a.db, id, data)
}

func (a *App) DeleteComputer(id int) error {
	if a.db == nil {
		return errNoDB
	}
	return services.DeleteComputer(a.db, id)
}

// ---------------------------------------------------------------------------
// Smartphone write methods
// ---------------------------------------------------------------------------

func (a *App) AddSmartphone(data models.SmartphoneInput) error {
	if a.db == nil {
		return errNoDB
	}
	_, err := services.InsertSmartphone(a.db, data)
	return err
}

func (a *App) UpdateSmartphone(id int, data models.SmartphoneInput) error {
	if a.db == nil {
		return errNoDB
	}
	return services.UpdateSmartphone(a.db, id, data)
}

func (a *App) DeleteSmartphone(id int) error {
	if a.db == nil {
		return errNoDB
	}
	return services.DeleteSmartphone(a.db, id)
}

// ---------------------------------------------------------------------------
// Tablet write methods
// ---------------------------------------------------------------------------

func (a *App) AddTablet(data models.TabletInput) error {
	if a.db == nil {
		return errNoDB
	}
	_, err := services.InsertTablet(a.db, data)
	return err
}

func (a *App) UpdateTablet(id int, data models.TabletInput) error {
	if a.db == nil {
		return errNoDB
	}
	return services.UpdateTablet(a.db, id, data)
}

func (a *App) DeleteTablet(id int) error {
	if a.db == nil {
		return errNoDB
	}
	return services.DeleteTablet(a.db, id)
}

// ---------------------------------------------------------------------------
// WindowsKey write methods
// ---------------------------------------------------------------------------

func (a *App) AddWindowsKey(data models.WindowsKeyInput) error {
	if a.db == nil {
		return errNoDB
	}
	_, err := services.InsertWindowsKey(a.db, data)
	return err
}

func (a *App) UpdateWindowsKey(id int, data models.WindowsKeyInput) error {
	if a.db == nil {
		return errNoDB
	}
	return services.UpdateWindowsKey(a.db, id, data)
}

func (a *App) DeleteWindowsKey(id int) error {
	if a.db == nil {
		return errNoDB
	}
	return services.DeleteWindowsKey(a.db, id)
}

// ---------------------------------------------------------------------------
// Antivirus write methods
// ---------------------------------------------------------------------------

func (a *App) AddAntivirus(data models.AntivirusInput) error {
	if a.db == nil {
		return errNoDB
	}
	_, err := services.InsertAntivirus(a.db, data)
	return err
}

func (a *App) UpdateAntivirus(id int, data models.AntivirusInput) error {
	if a.db == nil {
		return errNoDB
	}
	return services.UpdateAntivirus(a.db, id, data)
}

func (a *App) DeleteAntivirus(id int) error {
	if a.db == nil {
		return errNoDB
	}
	return services.DeleteAntivirus(a.db, id)
}

// ---------------------------------------------------------------------------
// OtherSoftware write methods
// ---------------------------------------------------------------------------

func (a *App) AddOtherSoftware(data models.OtherSoftwareInput) error {
	if a.db == nil {
		return errNoDB
	}
	_, err := services.InsertOtherSoftware(a.db, data)
	return err
}

func (a *App) UpdateOtherSoftware(id int, data models.OtherSoftwareInput) error {
	if a.db == nil {
		return errNoDB
	}
	return services.UpdateOtherSoftware(a.db, id, data)
}

func (a *App) DeleteOtherSoftware(id int) error {
	if a.db == nil {
		return errNoDB
	}
	return services.DeleteOtherSoftware(a.db, id)
}

// ---------------------------------------------------------------------------
// User write methods
// ---------------------------------------------------------------------------

func (a *App) AddUser(data models.UserInput) error {
	if a.db == nil {
		return errNoDB
	}
	_, err := services.InsertUser(a.db, data)
	return err
}

func (a *App) UpdateUser(id int, data models.UserInput) error {
	if a.db == nil {
		return errNoDB
	}
	return services.UpdateUser(a.db, id, data)
}

func (a *App) DeleteUser(id int) error {
	if a.db == nil {
		return errNoDB
	}
	return services.DeleteUser(a.db, id)
}

// ---------------------------------------------------------------------------
// Dropdown helpers
// ---------------------------------------------------------------------------

func (a *App) ListUsersForDropdown() ([]models.DropdownItem, error) {
	if a.db == nil {
		return nil, errNoDB
	}
	out, err := services.ListUsersForDropdown(a.db)
	if out == nil {
		out = []models.DropdownItem{}
	}
	return out, err
}

func (a *App) ListDevicesForDropdown() ([]models.DeviceDropdownItem, error) {
	if a.db == nil {
		return nil, errNoDB
	}
	out, err := services.ListDevicesForDropdown(a.db)
	if out == nil {
		out = []models.DeviceDropdownItem{}
	}
	return out, err
}

var _ = backup.BackupDB
var _ = database.Open // keeps modernc.org/sqlite in go.mod
