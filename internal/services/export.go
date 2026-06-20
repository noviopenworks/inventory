package services

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
)

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// BuildCSV returns header + data rows for a category in canonical column order.
//
//nolint:gocyclo // flat query assembly: each case is an independent category branch with no nesting
func BuildCSV(db *sql.DB, category string) ([][]string, error) {
	switch category {
	case "computers":
		items, err := ListComputers(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"Name", "Model", "Status", "Purchase Date", "Warranty Expiry"}}
		for _, c := range items {
			out = append(out, []string{c.Name, c.Model, c.Status, deref(c.PurchaseDate), deref(c.WarrantyExpiry)})
		}
		return out, nil
	case "smartphones":
		items, err := ListSmartphones(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"Name", "Model", "Status", "Purchase Date", "Warranty Expiry"}}
		for _, c := range items {
			out = append(out, []string{c.Name, c.Model, c.Status, deref(c.PurchaseDate), deref(c.WarrantyExpiry)})
		}
		return out, nil
	case "tablets":
		items, err := ListTablets(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"Name", "Model", "Status", "Purchase Date", "Warranty Expiry"}}
		for _, c := range items {
			out = append(out, []string{c.Name, c.Model, c.Status, deref(c.PurchaseDate), deref(c.WarrantyExpiry)})
		}
		return out, nil
	case "windowskeys":
		items, err := ListWindowsKeys(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"License Key", "Status"}}
		for _, w := range items {
			out = append(out, []string{w.LicenseKey, w.Status})
		}
		return out, nil
	case "antivirus":
		items, err := ListAntivirus(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"Name", "License Key", "Status", "Expiry Date"}}
		for _, a := range items {
			out = append(out, []string{a.Name, a.LicenseKey, a.Status, deref(a.ExpiryDate)})
		}
		return out, nil
	case "othersoftware":
		items, err := ListOtherSoftware(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"Name", "License Key", "Status", "Expiry Date"}}
		for _, s := range items {
			out = append(out, []string{s.Name, s.LicenseKey, s.Status, deref(s.ExpiryDate)})
		}
		return out, nil
	case "users":
		items, err := ListUsers(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"Name", "Surname", "Status"}}
		for _, u := range items {
			out = append(out, []string{u.Name, deref(u.Surname), u.Status})
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unknown category: %s", category)
	}
}

// WriteCSV writes rows to path using encoding/csv.
func WriteCSV(path string, rows [][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	w := csv.NewWriter(f)
	if err := w.WriteAll(rows); err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}
