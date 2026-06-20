package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
	"inventory/internal/models"
)

func TestListSmartphones_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListSmartphones()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAddSmartphone_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddSmartphone(models.SmartphoneInput{Name: "X", Model: "Y", Status: "active"})
	assert.Error(t, err)
}

func TestAddSmartphone_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddSmartphone(models.SmartphoneInput{Name: "Pixel", Model: "9 Pro", Status: "active"}))
	list, err := app.ListSmartphones()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Pixel", list[0].Name)
}

func TestUpdateSmartphone(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddSmartphone(models.SmartphoneInput{Name: "Old", Model: "M", Status: "active"}))
	list, _ := app.ListSmartphones()
	id := list[0].ID

	require.NoError(t, app.UpdateSmartphone(id, models.SmartphoneInput{Name: "New", Model: "M", Status: "active"}))
	list, _ = app.ListSmartphones()
	assert.Equal(t, "New", list[0].Name)
}

func TestDeleteSmartphone(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddSmartphone(models.SmartphoneInput{Name: "S", Model: "M", Status: "active"}))
	list, _ := app.ListSmartphones()
	require.NoError(t, app.DeleteSmartphone(list[0].ID))
	list, _ = app.ListSmartphones()
	assert.Empty(t, list)
}
