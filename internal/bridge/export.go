package bridge

import (
	"database/sql"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"inventory/internal/services"
)

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
