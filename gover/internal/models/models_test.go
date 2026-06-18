package models_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gover/internal/models"
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

func TestComputerInput_JSONRoundTrip(t *testing.T) {
	userID := 5
	pd := "2024-01-01"
	we := "2027-01-01"
	notes := "some notes"
	in := models.ComputerInput{
		Name:           "MyPC",
		Model:          "Dell XPS",
		UserID:         &userID,
		Status:         "active",
		PurchaseDate:   &pd,
		WarrantyExpiry: &we,
		Notes:          &notes,
	}
	data, err := json.Marshal(in)
	require.NoError(t, err)

	// Verify JSON keys are camelCase
	assert.Contains(t, string(data), `"userId"`)
	assert.Contains(t, string(data), `"purchaseDate"`)
	assert.Contains(t, string(data), `"warrantyExpiry"`)

	var got models.ComputerInput
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, in, got)
}

func TestUserInput_JSONRoundTrip(t *testing.T) {
	surname := "Smith"
	notes := "admin user"
	in := models.UserInput{
		Name:    "Alice",
		Surname: &surname,
		Status:  "active",
		Notes:   &notes,
	}
	data, err := json.Marshal(in)
	require.NoError(t, err)
	var got models.UserInput
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, in, got)
}

func TestDropdownItem_JSONRoundTrip(t *testing.T) {
	item := models.DropdownItem{ID: 3, Name: "Alice"}
	data, err := json.Marshal(item)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"id"`)
	assert.Contains(t, string(data), `"name"`)
	var got models.DropdownItem
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, item, got)
}

func TestDeviceDropdownItem_JSONRoundTrip(t *testing.T) {
	item := models.DeviceDropdownItem{ID: 7, Name: "PC1", Kind: "computer"}
	data, err := json.Marshal(item)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"kind"`)
	var got models.DeviceDropdownItem
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, item, got)
}
