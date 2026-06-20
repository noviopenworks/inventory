package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/services"
)

func TestGetAlerts_ExpiringWithin30Days(t *testing.T) {
	db := openTestDB(t)
	// expiry in 15 days — should appear
	expiring := time.Now().AddDate(0, 0, 15).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO antivirus (name, status, expiry_date) VALUES (?, ?, ?)`, "AVG", "Active", expiring)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "expiring", result[0].Severity)
}

func TestGetAlerts_NotExpiringSoon(t *testing.T) {
	db := openTestDB(t)
	// expiry in 45 days — should NOT appear with 30-day window
	notExpiring := time.Now().AddDate(0, 0, 45).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO antivirus (name, status, expiry_date) VALUES (?, ?, ?)`, "AVG", "Active", notExpiring)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Empty(t, result)

	// Also test other_software expiring soon
	expiringSoon := time.Now().AddDate(0, 0, 10).Format("2006-01-02")
	_, err = db.Exec(`INSERT INTO other_software (name, status, expiry_date) VALUES (?, ?, ?)`, "Slack", "Active", expiringSoon)
	require.NoError(t, err)
	result, err = services.GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetAlerts_AlreadyExpired(t *testing.T) {
	db := openTestDB(t)
	// expired yesterday — should appear with severity "expired"
	expired := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO antivirus (name, status, expiry_date) VALUES (?, ?, ?)`, "AVG", "Active", expired)
	require.NoError(t, err)
	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "expired", result[0].Severity)
}

func TestGetAlerts_ComputerWarrantyExpiring(t *testing.T) {
	db := openTestDB(t)
	warranty := time.Now().AddDate(0, 0, 15).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO computers (name, status, warranty_expiry) VALUES (?, ?, ?)`, "Dell XPS", "Active", warranty)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Computer", result[0].Category)
	assert.Equal(t, "expiring", result[0].Severity)
	assert.Equal(t, "Dell XPS", result[0].Name)
}

func TestGetAlerts_SmartphoneWarrantyExpired(t *testing.T) {
	db := openTestDB(t)
	warranty := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO smartphones (name, status, warranty_expiry) VALUES (?, ?, ?)`, "iPhone", "Active", warranty)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Smartphone", result[0].Category)
	assert.Equal(t, "expired", result[0].Severity)
}

func TestGetAlerts_TabletWarrantyCovered(t *testing.T) {
	db := openTestDB(t)
	warranty := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO tablets (name, status, warranty_expiry) VALUES (?, ?, ?)`, "iPad", "Active", warranty)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Tablet", result[0].Category)
	assert.Equal(t, "expiring", result[0].Severity)
}

func TestGetAlerts_RetiredExcluded(t *testing.T) {
	db := openTestDB(t)
	expired := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO antivirus (name, status, expiry_date) VALUES (?, ?, ?)`, "Old AV", "Retired", expired)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Empty(t, result)
}
