package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/config"
	"inventory/internal/models"
)

func TestLoad_Defaults(t *testing.T) {
	// Point XDG_CONFIG_HOME to a temp dir with no config file
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "comfortable", cfg.Density)
	assert.False(t, cfg.DarkMode)
	assert.Equal(t, 30, cfg.ExpiryWarningDays)
	// DBPath should be non-empty (default path)
	assert.NotEmpty(t, cfg.DBPath)
}

func TestSave_And_Load(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	want := models.AppConfig{
		DBPath:            "/some/path/inventory.db",
		Density:           "compact",
		DarkMode:          true,
		ExpiryWarningDays: 14,
	}
	require.NoError(t, config.Save(want))

	got, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestSave_CreatesDirectory(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	// The inventory/ subdirectory does not exist yet
	require.NoError(t, config.Save(models.AppConfig{DBPath: "/db/path"}))

	// File must exist
	p := filepath.Join(tmp, "inventory", "config.json")
	_, err := os.Stat(p)
	assert.NoError(t, err)

	// Must be valid JSON
	data, err := os.ReadFile(p)
	require.NoError(t, err)
	var m map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &m))
}
