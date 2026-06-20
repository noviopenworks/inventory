package bridge_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
	"inventory/internal/database"
)

func TestApp_NewDatabase(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	app := bridge.NewApp()
	path := filepath.Join(t.TempDir(), "new.db")
	err := app.NewDatabase(path)
	require.NoError(t, err)

	// Database should now be open: list should succeed
	list, err := app.ListComputers()
	require.NoError(t, err)
	assert.NotNil(t, list)
}

func TestApp_NewDatabase_BadPath(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	app := bridge.NewApp()
	err := app.NewDatabase("/nonexistent-dir/sub/new.db")
	require.Error(t, err)
}

func TestApp_OpenDatabase(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	app := bridge.NewApp()

	// Create a real db file first
	dbPath := filepath.Join(t.TempDir(), "existing.db")
	require.NoError(t, app.NewDatabase(dbPath))

	// Now open a fresh app and open the existing db
	app2 := bridge.NewApp()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	err := app2.OpenDatabase(dbPath)
	require.NoError(t, err)

	list, err := app2.ListComputers()
	require.NoError(t, err)
	assert.NotNil(t, list)
}

func TestApp_OpenDatabase_NotFound(t *testing.T) {
	app := bridge.NewApp()
	err := app.OpenDatabase("/nonexistent/path/db.db")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}

func TestApp_BackupDatabaseDialog_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.BackupDatabaseDialog()
	require.ErrorIs(t, err, bridge.ErrNoDB)
}

func TestApp_BackupDatabase_NoDB(t *testing.T) {
	app := bridge.NewApp()
	require.ErrorIs(t, app.BackupDatabase(filepath.Join(t.TempDir(), "x.db")), bridge.ErrNoDB)
}

func TestApp_BackupDatabase_WritesCopy(t *testing.T) {
	db, err := database.Open(filepath.Join(t.TempDir(), "live.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, database.InitSchema(db))
	_, err = db.Exec(`INSERT INTO users (name, surname) VALUES (?, ?)`, "Alice", "Smith")
	require.NoError(t, err)

	app := bridge.NewAppWithDB(db)
	dest := filepath.Join(t.TempDir(), "copy.db")
	require.NoError(t, app.BackupDatabase(dest))

	copyDB, err := database.Open(dest)
	require.NoError(t, err)
	t.Cleanup(func() { _ = copyDB.Close() })
	var count int
	require.NoError(t, copyDB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count))
	require.Equal(t, 1, count)
}
