package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
	"inventory/internal/models"
)

func TestListWindowsKeys_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListWindowsKeys()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAddWindowsKey_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddWindowsKey(models.WindowsKeyInput{LicenseKey: "KEY", Status: "active"})
	assert.Error(t, err)
}

func TestAddWindowsKey_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddWindowsKey(models.WindowsKeyInput{LicenseKey: "ABCD-1234", Status: "active"}))
	list, err := app.ListWindowsKeys()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "ABCD-1234", list[0].LicenseKey)
}

func TestUpdateWindowsKey(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddWindowsKey(models.WindowsKeyInput{LicenseKey: "OLD-KEY", Status: "active"}))
	list, _ := app.ListWindowsKeys()
	id := list[0].ID

	require.NoError(t, app.UpdateWindowsKey(id, models.WindowsKeyInput{LicenseKey: "NEW-KEY", Status: "active"}))
	list, _ = app.ListWindowsKeys()
	assert.Equal(t, "NEW-KEY", list[0].LicenseKey)
}

func TestDeleteWindowsKey(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddWindowsKey(models.WindowsKeyInput{LicenseKey: "KEY", Status: "active"}))
	list, _ := app.ListWindowsKeys()
	require.NoError(t, app.DeleteWindowsKey(list[0].ID))
	list, _ = app.ListWindowsKeys()
	assert.Empty(t, list)
}
