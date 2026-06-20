package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
	"inventory/internal/models"
)

func TestListTablets_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListTablets()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAddTablet_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddTablet(models.TabletInput{Name: "X", Model: "Y", Status: "active"})
	assert.Error(t, err)
}

func TestAddTablet_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddTablet(models.TabletInput{Name: "iPad", Model: "Pro 13", Status: "active"}))
	list, err := app.ListTablets()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "iPad", list[0].Name)
}

func TestUpdateTablet(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddTablet(models.TabletInput{Name: "Old", Model: "M", Status: "active"}))
	list, _ := app.ListTablets()
	id := list[0].ID

	require.NoError(t, app.UpdateTablet(id, models.TabletInput{Name: "New", Model: "M", Status: "active"}))
	list, _ = app.ListTablets()
	assert.Equal(t, "New", list[0].Name)
}

func TestDeleteTablet(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddTablet(models.TabletInput{Name: "T", Model: "M", Status: "active"}))
	list, _ := app.ListTablets()
	require.NoError(t, app.DeleteTablet(list[0].ID))
	list, _ = app.ListTablets()
	assert.Empty(t, list)
}
