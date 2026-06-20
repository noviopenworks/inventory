package services

import (
	"database/sql"

	"inventory/internal/models"
)

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

// InsertOtherSoftware inserts a new other-software record and returns its id.
func InsertOtherSoftware(db *sql.DB, d models.OtherSoftwareInput) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO other_software (name, license_key, computer_id, smartphone_id, tablet_id, status, expiry_date, notes)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		d.Name, d.LicenseKey, d.ComputerID, d.SmartphoneID, d.TabletID, d.Status, d.ExpiryDate, d.Notes,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateOtherSoftware updates the other-software record identified by id.
func UpdateOtherSoftware(db *sql.DB, id int, d models.OtherSoftwareInput) error {
	_, err := db.Exec(
		`UPDATE other_software SET name=?, license_key=?, computer_id=?, smartphone_id=?, tablet_id=?,
		 status=?, expiry_date=?, notes=?, updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.LicenseKey, d.ComputerID, d.SmartphoneID, d.TabletID, d.Status, d.ExpiryDate, d.Notes, id,
	)
	return err
}

// DeleteOtherSoftware removes the other-software record identified by id.
func DeleteOtherSoftware(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM other_software WHERE id=?`, id)
	return err
}
