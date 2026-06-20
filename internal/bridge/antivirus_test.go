package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
	"inventory/internal/models"
)

func TestListAntivirus_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListAntivirus()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAddAntivirus_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddAntivirus(models.AntivirusInput{Name: "AV", LicenseKey: "KEY", Status: "active"})
	assert.Error(t, err)
}

func TestAddAntivirus_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddAntivirus(models.AntivirusInput{Name: "Norton", LicenseKey: "NRT-001", Status: "active"}))
	list, err := app.ListAntivirus()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Norton", list[0].Name)
}

func TestUpdateAntivirus(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddAntivirus(models.AntivirusInput{Name: "Old", LicenseKey: "K", Status: "active"}))
	list, _ := app.ListAntivirus()
	id := list[0].ID

	require.NoError(t, app.UpdateAntivirus(id, models.AntivirusInput{Name: "New", LicenseKey: "K", Status: "active"}))
	list, _ = app.ListAntivirus()
	assert.Equal(t, "New", list[0].Name)
}

func TestDeleteAntivirus(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddAntivirus(models.AntivirusInput{Name: "AV", LicenseKey: "K", Status: "active"}))
	list, _ := app.ListAntivirus()
	require.NoError(t, app.DeleteAntivirus(list[0].ID))
	list, _ = app.ListAntivirus()
	assert.Empty(t, list)
}
