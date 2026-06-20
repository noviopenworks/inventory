package services

import (
	"database/sql"

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

// InsertComputer inserts a new computer record and returns its id.
func InsertComputer(db *sql.DB, d models.ComputerInput) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO computers (name, model, user_id, status, purchase_date, warranty_expiry, notes)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateComputer updates the computer identified by id.
func UpdateComputer(db *sql.DB, id int, d models.ComputerInput) error {
	_, err := db.Exec(
		`UPDATE computers SET name=?, model=?, user_id=?, status=?, purchase_date=?,
		 warranty_expiry=?, notes=?, updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes, id,
	)
	return err
}

// DeleteComputer removes the computer identified by id.
func DeleteComputer(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM computers WHERE id=?`, id)
	return err
}
