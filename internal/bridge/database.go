package bridge

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"inventory/internal/backup"
	"inventory/internal/config"
	"inventory/internal/database"
	"inventory/internal/models"
)

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

// BackupDatabase writes a standalone copy of the open database to destPath.
func (a *App) BackupDatabase(destPath string) error {
	return a.withDB(func(db *sql.DB) error { return backup.Create(db, destPath) })
}

// BackupDatabaseDialog prompts for a destination and backs up the open
// database there. Returns nil if the user cancels.
func (a *App) BackupDatabaseDialog() error {
	// Guard before opening the dialog so we don't prompt when no DB is loaded.
	if a.db == nil {
		return ErrNoDB
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "inventory-backup-" + time.Now().Format("2006-01-02") + ".db",
		Filters:         []runtime.FileFilter{{DisplayName: "SQLite database", Pattern: "*.db"}},
	})
	if err != nil {
		return err
	}
	if path == "" {
		return nil // user cancelled
	}
	return a.BackupDatabase(path)
}
