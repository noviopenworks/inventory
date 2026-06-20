package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"inventory/internal/models"
)

// ListComputers returns all computers ordered by name.
func ListComputers(db *sql.DB) ([]models.Computer, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), COALESCE(model,''), user_id, COALESCE(status,''),
		        purchase_date, warranty_expiry, notes, created_at, updated_at
		 FROM computers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []models.Computer
	for rows.Next() {
		var c models.Computer
		if err := rows.Scan(&c.ID, &c.Name, &c.Model, &c.UserID, &c.Status,
			&c.PurchaseDate, &c.WarrantyExpiry, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// ListSmartphones returns all smartphones ordered by name.
func ListSmartphones(db *sql.DB) ([]models.Smartphone, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), COALESCE(model,''), user_id, COALESCE(status,''),
		        purchase_date, warranty_expiry, notes, created_at, updated_at
		 FROM smartphones ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []models.Smartphone
	for rows.Next() {
		var s models.Smartphone
		if err := rows.Scan(&s.ID, &s.Name, &s.Model, &s.UserID, &s.Status,
			&s.PurchaseDate, &s.WarrantyExpiry, &s.Notes, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// ListTablets returns all tablets ordered by name.
func ListTablets(db *sql.DB) ([]models.Tablet, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), COALESCE(model,''), user_id, COALESCE(status,''),
		        purchase_date, warranty_expiry, notes, created_at, updated_at
		 FROM tablets ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []models.Tablet
	for rows.Next() {
		var tab models.Tablet
		if err := rows.Scan(&tab.ID, &tab.Name, &tab.Model, &tab.UserID, &tab.Status,
			&tab.PurchaseDate, &tab.WarrantyExpiry, &tab.Notes, &tab.CreatedAt, &tab.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, tab)
	}
	return out, rows.Err()
}

// ListWindowsKeys returns all Windows license keys ordered by id.
func ListWindowsKeys(db *sql.DB) ([]models.WindowsKey, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(license_key,''), computer_id, COALESCE(status,''),
		        notes, created_at, updated_at
		 FROM windows_keys ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []models.WindowsKey
	for rows.Next() {
		var w models.WindowsKey
		if err := rows.Scan(&w.ID, &w.LicenseKey, &w.ComputerID, &w.Status,
			&w.Notes, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

// ListAntivirus returns all antivirus records ordered by name.
func ListAntivirus(db *sql.DB) ([]models.Antivirus, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), COALESCE(license_key,''), computer_id, smartphone_id, tablet_id,
		        COALESCE(status,''), expiry_date, notes, created_at, updated_at
		 FROM antivirus ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []models.Antivirus
	for rows.Next() {
		var a models.Antivirus
		if err := rows.Scan(&a.ID, &a.Name, &a.LicenseKey, &a.ComputerID, &a.SmartphoneID,
			&a.TabletID, &a.Status, &a.ExpiryDate, &a.Notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// ListOtherSoftware returns all other-software records ordered by name.
func ListOtherSoftware(db *sql.DB) ([]models.OtherSoftware, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), COALESCE(license_key,''), computer_id, smartphone_id, tablet_id,
		        COALESCE(status,''), expiry_date, notes, created_at, updated_at
		 FROM other_software ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []models.OtherSoftware
	for rows.Next() {
		var s models.OtherSoftware
		if err := rows.Scan(&s.ID, &s.Name, &s.LicenseKey, &s.ComputerID, &s.SmartphoneID,
			&s.TabletID, &s.Status, &s.ExpiryDate, &s.Notes, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// ListUsers returns all users ordered by name.
func ListUsers(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), surname, COALESCE(status,''), notes, created_at, updated_at
		 FROM users ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Status,
			&u.Notes, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// alertSources lists every (table, date column) pair that contributes to
// expiry/warranty alerts. Mirrors db/alerts.py: Windows Key is intentionally
// omitted (no expiry column); Users is irrelevant.
type alertSource struct {
	table    string
	category string
	dateCol  string
}

var alertSources = []alertSource{
	{"computers", "Computer", "warranty_expiry"},
	{"smartphones", "Smartphone", "warranty_expiry"},
	{"tablets", "Tablet", "warranty_expiry"},
	{"antivirus", "Antivirus", "expiry_date"},
	{"other_software", "Other Software", "expiry_date"},
}

// GetAlerts returns expiry and warranty alerts for items expiring within warningDays.
func GetAlerts(db *sql.DB, warningDays int) ([]models.Alert, error) {
	threshold := fmt.Sprintf("+%d days", warningDays)
	parts := make([]string, 0, len(alertSources))
	params := make([]any, 0, len(alertSources))
	for _, src := range alertSources {
		parts = append(parts, fmt.Sprintf(`
		SELECT '%s' AS category, id, COALESCE(name, '') AS name, %s AS expiry_date,
		       CAST(julianday(%s) - julianday('now') AS INTEGER) AS days_remaining
		FROM %s
		WHERE %s IS NOT NULL
		  AND date(%s) <= date('now', ?)
		  AND status != 'Retired'`,
			src.category, src.dateCol, src.dateCol, src.table, src.dateCol, src.dateCol))
		params = append(params, threshold)
	}
	query := strings.Join(parts, "\n UNION ALL \n") + "\n ORDER BY days_remaining ASC"

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	today := time.Now().Format("2006-01-02")
	var out []models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.Category, &a.ID, &a.Name, &a.ExpiryDate, &a.DaysRemaining); err != nil {
			return nil, err
		}
		if a.ExpiryDate < today {
			a.Severity = "expired"
		} else {
			a.Severity = "expiring"
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
