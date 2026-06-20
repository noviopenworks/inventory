package services

import (
	"database/sql"

	"inventory/internal/models"
)

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

// InsertAntivirus inserts a new antivirus record and returns its id.
func InsertAntivirus(db *sql.DB, d models.AntivirusInput) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO antivirus (name, license_key, computer_id, smartphone_id, tablet_id, status, expiry_date, notes)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		d.Name, d.LicenseKey, d.ComputerID, d.SmartphoneID, d.TabletID, d.Status, d.ExpiryDate, d.Notes,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateAntivirus updates the antivirus record identified by id.
func UpdateAntivirus(db *sql.DB, id int, d models.AntivirusInput) error {
	_, err := db.Exec(
		`UPDATE antivirus SET name=?, license_key=?, computer_id=?, smartphone_id=?, tablet_id=?,
		 status=?, expiry_date=?, notes=?, updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.LicenseKey, d.ComputerID, d.SmartphoneID, d.TabletID, d.Status, d.ExpiryDate, d.Notes, id,
	)
	return err
}

// DeleteAntivirus removes the antivirus record identified by id.
func DeleteAntivirus(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM antivirus WHERE id=?`, id)
	return err
}
