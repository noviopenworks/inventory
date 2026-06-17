package bridge

import (
	"gover2/internal/backup"
	"gover2/internal/config"
	"gover2/internal/database"
	"gover2/internal/models"
	"gover2/internal/services"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) ListComputers() ([]models.Computer, error) {
	return services.ListComputers(nil)
}

func (a *App) ListSmartphones() ([]models.Smartphone, error) {
	return services.ListSmartphones(nil)
}

func (a *App) ListTablets() ([]models.Tablet, error) {
	return services.ListTablets(nil)
}

func (a *App) ListWindowsKeys() ([]models.WindowsKey, error) {
	return services.ListWindowsKeys(nil)
}

func (a *App) ListAntivirus() ([]models.Antivirus, error) {
	return services.ListAntivirus(nil)
}

func (a *App) ListOtherSoftware() ([]models.OtherSoftware, error) {
	return services.ListOtherSoftware(nil)
}

func (a *App) ListUsers() ([]models.User, error) {
	return services.ListUsers(nil)
}

func (a *App) GetDatabasePath() (string, error) {
	return "", nil
}

func (a *App) GetConfig() (models.AppConfig, error) {
	return config.Load()
}

func (a *App) SetConfig(cfg models.AppConfig) error {
	return config.Save(cfg)
}

func (a *App) NewDatabase(path string) error {
	return nil
}

func (a *App) OpenDatabase(path string) error {
	return nil
}

func (a *App) ExportCSV(category string, destPath string) error {
	return nil
}

var _ = backup.BackupDB
var _ = database.Open
