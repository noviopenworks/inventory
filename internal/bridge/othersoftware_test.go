package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
	"inventory/internal/models"
)

func TestListOtherSoftware_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListOtherSoftware()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAddOtherSoftware_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddOtherSoftware(models.OtherSoftwareInput{Name: "SW", LicenseKey: "KEY", Status: "active"})
	assert.Error(t, err)
}

func TestAddOtherSoftware_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddOtherSoftware(models.OtherSoftwareInput{Name: "Office", LicenseKey: "OFF-001", Status: "active"}))
	list, err := app.ListOtherSoftware()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Office", list[0].Name)
}

func TestUpdateOtherSoftware(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddOtherSoftware(models.OtherSoftwareInput{Name: "Old", LicenseKey: "K", Status: "active"}))
	list, _ := app.ListOtherSoftware()
	id := list[0].ID

	require.NoError(t, app.UpdateOtherSoftware(id, models.OtherSoftwareInput{Name: "New", LicenseKey: "K", Status: "active"}))
	list, _ = app.ListOtherSoftware()
	assert.Equal(t, "New", list[0].Name)
}

func TestDeleteOtherSoftware(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddOtherSoftware(models.OtherSoftwareInput{Name: "SW", LicenseKey: "K", Status: "active"}))
	list, _ := app.ListOtherSoftware()
	require.NoError(t, app.DeleteOtherSoftware(list[0].ID))
	list, _ = app.ListOtherSoftware()
	assert.Empty(t, list)
}
