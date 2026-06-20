package services

import (
	"database/sql"

	"inventory/internal/models"
)

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

// InsertWindowsKey inserts a new Windows license key record and returns its id.
func InsertWindowsKey(db *sql.DB, d models.WindowsKeyInput) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO windows_keys (license_key, computer_id, status, notes)
		 VALUES (?, ?, ?, ?)`,
		d.LicenseKey, d.ComputerID, d.Status, d.Notes,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateWindowsKey updates the Windows license key identified by id.
func UpdateWindowsKey(db *sql.DB, id int, d models.WindowsKeyInput) error {
	_, err := db.Exec(
		`UPDATE windows_keys SET license_key=?, computer_id=?, status=?, notes=?,
		 updated_at=datetime('now','localtime') WHERE id=?`,
		d.LicenseKey, d.ComputerID, d.Status, d.Notes, id,
	)
	return err
}

// DeleteWindowsKey removes the Windows license key identified by id.
func DeleteWindowsKey(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM windows_keys WHERE id=?`, id)
	return err
}
