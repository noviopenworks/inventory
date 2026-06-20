// Package backup creates standalone copies of the live SQLite database.
package backup

import (
	"database/sql"
	"fmt"
	"strings"
)

// Create writes a consistent standalone copy of db to destPath using SQLite's
// VACUUM INTO, which produces a defragmented snapshot regardless of WAL state.
// destPath must not already exist (VACUUM INTO refuses to overwrite).
func Create(db *sql.DB, destPath string) error {
	// VACUUM INTO does not accept a bound parameter for its destination, so the
	// path is inlined with single quotes escaped by doubling.
	escaped := strings.ReplaceAll(destPath, "'", "''")
	if _, err := db.Exec(fmt.Sprintf("VACUUM INTO '%s'", escaped)); err != nil {
		return fmt.Errorf("backup database: %w", err)
	}
	return nil
}
