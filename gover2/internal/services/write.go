package services

import (
	"database/sql"

	"gover2/internal/models"
)

// ---------------------------------------------------------------------------
// Computers
// ---------------------------------------------------------------------------

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

func UpdateComputer(db *sql.DB, id int, d models.ComputerInput) error {
	_, err := db.Exec(
		`UPDATE computers SET name=?, model=?, user_id=?, status=?, purchase_date=?,
		 warranty_expiry=?, notes=?, updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes, id,
	)
	return err
}

func DeleteComputer(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM computers WHERE id=?`, id)
	return err
}

// ---------------------------------------------------------------------------
// Smartphones
// ---------------------------------------------------------------------------

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

func UpdateSmartphone(db *sql.DB, id int, d models.SmartphoneInput) error {
	_, err := db.Exec(
		`UPDATE smartphones SET name=?, model=?, user_id=?, status=?, purchase_date=?,
		 warranty_expiry=?, notes=?, updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes, id,
	)
	return err
}

func DeleteSmartphone(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM smartphones WHERE id=?`, id)
	return err
}

// ---------------------------------------------------------------------------
// Tablets
// ---------------------------------------------------------------------------

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

func UpdateTablet(db *sql.DB, id int, d models.TabletInput) error {
	_, err := db.Exec(
		`UPDATE tablets SET name=?, model=?, user_id=?, status=?, purchase_date=?,
		 warranty_expiry=?, notes=?, updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes, id,
	)
	return err
}

func DeleteTablet(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM tablets WHERE id=?`, id)
	return err
}

// ---------------------------------------------------------------------------
// Windows Keys
// ---------------------------------------------------------------------------

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

func UpdateWindowsKey(db *sql.DB, id int, d models.WindowsKeyInput) error {
	_, err := db.Exec(
		`UPDATE windows_keys SET license_key=?, computer_id=?, status=?, notes=?,
		 updated_at=datetime('now','localtime') WHERE id=?`,
		d.LicenseKey, d.ComputerID, d.Status, d.Notes, id,
	)
	return err
}

func DeleteWindowsKey(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM windows_keys WHERE id=?`, id)
	return err
}

// ---------------------------------------------------------------------------
// Antivirus
// ---------------------------------------------------------------------------

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

func UpdateAntivirus(db *sql.DB, id int, d models.AntivirusInput) error {
	_, err := db.Exec(
		`UPDATE antivirus SET name=?, license_key=?, computer_id=?, smartphone_id=?, tablet_id=?,
		 status=?, expiry_date=?, notes=?, updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.LicenseKey, d.ComputerID, d.SmartphoneID, d.TabletID, d.Status, d.ExpiryDate, d.Notes, id,
	)
	return err
}

func DeleteAntivirus(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM antivirus WHERE id=?`, id)
	return err
}

// ---------------------------------------------------------------------------
// OtherSoftware
// ---------------------------------------------------------------------------

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

func UpdateOtherSoftware(db *sql.DB, id int, d models.OtherSoftwareInput) error {
	_, err := db.Exec(
		`UPDATE other_software SET name=?, license_key=?, computer_id=?, smartphone_id=?, tablet_id=?,
		 status=?, expiry_date=?, notes=?, updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.LicenseKey, d.ComputerID, d.SmartphoneID, d.TabletID, d.Status, d.ExpiryDate, d.Notes, id,
	)
	return err
}

func DeleteOtherSoftware(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM other_software WHERE id=?`, id)
	return err
}

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

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

func UpdateUser(db *sql.DB, id int, d models.UserInput) error {
	_, err := db.Exec(
		`UPDATE users SET name=?, surname=?, status=?, notes=?,
		 updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.Surname, d.Status, d.Notes, id,
	)
	return err
}

func DeleteUser(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM users WHERE id=?`, id)
	return err
}
