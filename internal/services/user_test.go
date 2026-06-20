package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/models"
	"inventory/internal/services"
)

func TestListUsers_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO users (name, surname) VALUES (?, ?)`, "Alice", "Smith")
	require.NoError(t, err)
	result, err := services.ListUsers(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Alice", result[0].Name)
}

func TestUser_InsertUpdateDelete(t *testing.T) {
	db := openTestDB(t)

	id, err := services.InsertUser(db, models.UserInput{
		Name:   "Alice",
		Status: "active",
	})
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	list, err := services.ListUsers(db)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Alice", list[0].Name)

	err = services.UpdateUser(db, int(id), models.UserInput{
		Name:   "Bob",
		Status: "inactive",
	})
	require.NoError(t, err)

	list, err = services.ListUsers(db)
	require.NoError(t, err)
	assert.Equal(t, "Bob", list[0].Name)

	err = services.DeleteUser(db, int(id))
	require.NoError(t, err)

	list, err = services.ListUsers(db)
	require.NoError(t, err)
	assert.Empty(t, list)
}
