package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gover2/internal/bridge"
	"gover2/internal/database"
)

func newBridgeWithDB(t *testing.T) *bridge.App {
	t.Helper()
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, database.InitSchema(db))
	t.Cleanup(func() { db.Close() })
	return bridge.NewAppWithDB(db)
}

func TestListComputers_NoDB(t *testing.T) {
	app := bridge.NewApp()
	_, err := app.ListComputers()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no database")
}

func TestListComputers_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListComputers()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListSmartphones_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListSmartphones()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListTablets_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListTablets()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListWindowsKeys_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListWindowsKeys()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListAntivirus_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListAntivirus()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListOtherSoftware_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListOtherSoftware()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListUsers_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListUsers()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetDatabasePath_NoConfig(t *testing.T) {
	app := bridge.NewApp()
	path, err := app.GetDatabasePath()
	require.NoError(t, err)
	assert.Equal(t, "", path)
}

func TestGetConfig_Default(t *testing.T) {
	app := bridge.NewApp()
	cfg, err := app.GetConfig()
	require.NoError(t, err)
	assert.Equal(t, "comfortable", cfg.Density)
}
