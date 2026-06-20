package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
	"inventory/internal/models"
)

func TestListUsersForDropdown_NoDB(t *testing.T) {
	app := bridge.NewApp()
	_, err := app.ListUsersForDropdown()
	assert.Error(t, err)
}

func TestListDevicesForDropdown_NoDB(t *testing.T) {
	app := bridge.NewApp()
	_, err := app.ListDevicesForDropdown()
	assert.Error(t, err)
}

func TestListUsersForDropdown(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddUser(models.UserInput{Name: "Alice", Status: "active"}))
	require.NoError(t, app.AddUser(models.UserInput{Name: "Bob", Status: "active"}))

	items, err := app.ListUsersForDropdown()
	require.NoError(t, err)
	require.Len(t, items, 2)
	// Ordered by name: Alice, Bob
	assert.Equal(t, "Alice", items[0].Name)
	assert.Equal(t, "Bob", items[1].Name)
}

func TestListDevicesForDropdown(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddComputer(models.ComputerInput{Name: "PC1", Model: "M", Status: "active"}))
	require.NoError(t, app.AddSmartphone(models.SmartphoneInput{Name: "Phone1", Model: "M", Status: "active"}))
	require.NoError(t, app.AddTablet(models.TabletInput{Name: "Tablet1", Model: "M", Status: "active"}))

	items, err := app.ListDevicesForDropdown()
	require.NoError(t, err)
	require.Len(t, items, 3)

	// Verify all three kinds are present
	kinds := make(map[string]bool)
	for _, item := range items {
		kinds[item.Kind] = true
	}
	assert.True(t, kinds["computer"])
	assert.True(t, kinds["smartphone"])
	assert.True(t, kinds["tablet"])
}
