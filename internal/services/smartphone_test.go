package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/models"
	"inventory/internal/services"
)

func TestListSmartphones_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO smartphones (name, model, status) VALUES (?, ?, ?)`, "iPhone", "15 Pro", "Active")
	require.NoError(t, err)
	result, err := services.ListSmartphones(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "iPhone", result[0].Name)
}

func TestSmartphone_InsertUpdateDelete(t *testing.T) {
	db := openTestDB(t)

	id, err := services.InsertSmartphone(db, models.SmartphoneInput{
		Name:   "iPhone",
		Model:  "15 Pro",
		Status: "active",
	})
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	list, err := services.ListSmartphones(db)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "iPhone", list[0].Name)

	err = services.UpdateSmartphone(db, int(id), models.SmartphoneInput{
		Name:   "Samsung",
		Model:  "Galaxy S24",
		Status: "inactive",
	})
	require.NoError(t, err)

	list, err = services.ListSmartphones(db)
	require.NoError(t, err)
	assert.Equal(t, "Samsung", list[0].Name)

	err = services.DeleteSmartphone(db, int(id))
	require.NoError(t, err)

	list, err = services.ListSmartphones(db)
	require.NoError(t, err)
	assert.Empty(t, list)
}
