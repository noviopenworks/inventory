package services

import (
	"database/sql"

	"inventory/internal/models"
)

// ---------------------------------------------------------------------------
// Computers
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Smartphones
// ---------------------------------------------------------------------------

// InsertSmartphone inserts a new smartphone record and returns its id.
func InsertSmartphone(db *sql.DB, d models.SmartphoneInput) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO smartphones (name, model, user_id, status, purchase_date, warranty_expiry, notes)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateSmartphone updates the smartphone identified by id.
func UpdateSmartphone(db *sql.DB, id int, d models.SmartphoneInput) error {
	_, err := db.Exec(
		`UPDATE smartphones SET name=?, model=?, user_id=?, status=?, purchase_date=?,
		 warranty_expiry=?, notes=?, updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes, id,
	)
	return err
}

// DeleteSmartphone removes the smartphone identified by id.
func DeleteSmartphone(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM smartphones WHERE id=?`, id)
	return err
}

// ---------------------------------------------------------------------------
// Tablets
// ---------------------------------------------------------------------------

// InsertTablet inserts a new tablet record and returns its id.
func InsertTablet(db *sql.DB, d models.TabletInput) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO tablets (name, model, user_id, status, purchase_date, warranty_expiry, notes)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateTablet updates the tablet identified by id.
func UpdateTablet(db *sql.DB, id int, d models.TabletInput) error {
	_, err := db.Exec(
		`UPDATE tablets SET name=?, model=?, user_id=?, status=?, purchase_date=?,
		 warranty_expiry=?, notes=?, updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes, id,
	)
	return err
}

// DeleteTablet removes the tablet identified by id.
func DeleteTablet(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM tablets WHERE id=?`, id)
	return err
}

// ---------------------------------------------------------------------------
// Windows Keys
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Antivirus
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// OtherSoftware
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

// InsertUser inserts a new user record and returns its id.
func InsertUser(db *sql.DB, d models.UserInput) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO users (name, surname, status, notes)
		 VALUES (?, ?, ?, ?)`,
		d.Name, d.Surname, d.Status, d.Notes,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateUser updates the user identified by id.
func UpdateUser(db *sql.DB, id int, d models.UserInput) error {
	_, err := db.Exec(
		`UPDATE users SET name=?, surname=?, status=?, notes=?,
		 updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.Surname, d.Status, d.Notes, id,
	)
	return err
}

// DeleteUser removes the user identified by id.
func DeleteUser(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM users WHERE id=?`, id)
	return err
}
