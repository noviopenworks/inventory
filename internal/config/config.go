package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"inventory/internal/models"
)

func configPath() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "inventory", "config.json")
}

func defaultDBPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "inventory", "inventory.db")
}

func Load() (models.AppConfig, error) {
	defaults := models.AppConfig{
		DBPath:            defaultDBPath(),
		Density:           "comfortable",
		DarkMode:          false,
		ExpiryWarningDays: 30,
	}
	data, err := os.ReadFile(configPath())
	if os.IsNotExist(err) {
		return defaults, nil
	}
	if err != nil {
		return defaults, err
	}
	var cfg models.AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaults, err
	}
	// Fill zero-value fields with defaults
	if cfg.Density == "" {
		cfg.Density = defaults.Density
	}
	if cfg.ExpiryWarningDays == 0 {
		cfg.ExpiryWarningDays = defaults.ExpiryWarningDays
	}
	if cfg.DBPath == "" {
		cfg.DBPath = defaults.DBPath
	}
	return cfg, nil
}

func Save(cfg models.AppConfig) error {
	p := configPath()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	// Atomic write: write to .tmp then rename
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}
