package bridge

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"gover2/internal/backup"
	"gover2/internal/config"
	"gover2/internal/database"
	"gover2/internal/models"
	"gover2/internal/services"
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
		db.Close()
		return err
	}
	if a.db != nil {
		a.db.Close()
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

func (a *App) ExportCSV(category string, destPath string) error {
	return nil // Phase 5
}

var _ = backup.BackupDB
var _ = database.Open // keeps modernc.org/sqlite in go.mod
