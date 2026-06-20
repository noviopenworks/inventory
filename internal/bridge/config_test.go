package bridge_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
	"inventory/internal/models"
)

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

func TestApp_GetSetConfig_RoundTrip(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	app := bridge.NewApp()
	cfg := models.AppConfig{Density: "compact", DarkMode: true, ExpiryWarningDays: 7, DBPath: filepath.Join(t.TempDir(), "test.db")}
	require.NoError(t, app.SetConfig(cfg))
	got, err := app.GetConfig()
	require.NoError(t, err)
	require.Equal(t, cfg, got)
}
