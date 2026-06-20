package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
	"inventory/internal/models"
)

func TestListUsers_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListUsers()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAddUser_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddUser(models.UserInput{Name: "Alice", Status: "active"})
	assert.Error(t, err)
}

func TestAddUser_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddUser(models.UserInput{Name: "Alice", Status: "active"}))
	list, err := app.ListUsers()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Alice", list[0].Name)
}

func TestUpdateUser(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddUser(models.UserInput{Name: "Old", Status: "active"}))
	list, _ := app.ListUsers()
	id := list[0].ID

	require.NoError(t, app.UpdateUser(id, models.UserInput{Name: "New", Status: "active"}))
	list, _ = app.ListUsers()
	assert.Equal(t, "New", list[0].Name)
}

func TestDeleteUser(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddUser(models.UserInput{Name: "Bob", Status: "active"}))
	list, _ := app.ListUsers()
	require.NoError(t, app.DeleteUser(list[0].ID))
	list, _ = app.ListUsers()
	assert.Empty(t, list)
}
