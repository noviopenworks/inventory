package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
	"inventory/internal/models"
)

func TestListComputers_NoDB(t *testing.T) {
	app := bridge.NewApp()
	_, err := app.ListComputers()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no database")
}

func TestListComputers_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListComputers()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAddComputer_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddComputer(models.ComputerInput{Name: "X", Model: "Y", Status: "active"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no database")
}

func TestAddComputer_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	input := models.ComputerInput{Name: "ThinkPad", Model: "X1 Carbon", Status: "active"}
	err := app.AddComputer(input)
	require.NoError(t, err)

	list, err := app.ListComputers()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "ThinkPad", list[0].Name)
	assert.Equal(t, "X1 Carbon", list[0].Model)
}

func TestUpdateComputer(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddComputer(models.ComputerInput{Name: "OldName", Model: "M1", Status: "active"}))
	list, err := app.ListComputers()
	require.NoError(t, err)
	require.Len(t, list, 1)
	id := list[0].ID

	require.NoError(t, app.UpdateComputer(id, models.ComputerInput{Name: "NewName", Model: "M1", Status: "active"}))

	list, err = app.ListComputers()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "NewName", list[0].Name)
}

func TestDeleteComputer(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddComputer(models.ComputerInput{Name: "ToDelete", Model: "M2", Status: "active"}))
	list, err := app.ListComputers()
	require.NoError(t, err)
	require.Len(t, list, 1)
	id := list[0].ID

	require.NoError(t, app.DeleteComputer(id))

	list, err = app.ListComputers()
	require.NoError(t, err)
	assert.Empty(t, list)
}
