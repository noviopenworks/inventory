package config

import "gover2/internal/models"

func Load() (models.AppConfig, error) {
	return models.AppConfig{
		Density:           "comfortable",
		DarkMode:          false,
		ExpiryWarningDays: 30,
	}, nil
}

func Save(cfg models.AppConfig) error {
	return nil
}
