package models_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gover2/internal/models"
)

func TestComputer_JSONRoundTrip(t *testing.T) {
	userID := 1
	purchaseDate := "2023-01-15"
	warrantyExpiry := "2026-01-15"
	notes := "test notes"
	c := models.Computer{
		ID: 1, Name: "TestPC", Model: "Dell XPS", UserID: &userID,
		Status: "Active", PurchaseDate: &purchaseDate, WarrantyExpiry: &warrantyExpiry,
		Notes: &notes, CreatedAt: "2023-01-01", UpdatedAt: "2023-01-01",
	}
	data, err := json.Marshal(c)
	require.NoError(t, err)
	var got models.Computer
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, c, got)
}

func TestAppConfig_JSONRoundTrip(t *testing.T) {
	cfg := models.AppConfig{DBPath: "/db/path", Density: "compact", DarkMode: true, ExpiryWarningDays: 14}
	data, err := json.Marshal(cfg)
	require.NoError(t, err)
	var got models.AppConfig
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, cfg, got)
}
