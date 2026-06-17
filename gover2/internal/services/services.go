package services

import (
	"database/sql"

	"gover2/internal/models"
)

func ListComputers(db *sql.DB) ([]models.Computer, error) {
	return []models.Computer{}, nil
}

func ListSmartphones(db *sql.DB) ([]models.Smartphone, error) {
	return []models.Smartphone{}, nil
}

func ListTablets(db *sql.DB) ([]models.Tablet, error) {
	return []models.Tablet{}, nil
}

func ListWindowsKeys(db *sql.DB) ([]models.WindowsKey, error) {
	return []models.WindowsKey{}, nil
}

func ListAntivirus(db *sql.DB) ([]models.Antivirus, error) {
	return []models.Antivirus{}, nil
}

func ListOtherSoftware(db *sql.DB) ([]models.OtherSoftware, error) {
	return []models.OtherSoftware{}, nil
}

func ListUsers(db *sql.DB) ([]models.User, error) {
	return []models.User{}, nil
}
