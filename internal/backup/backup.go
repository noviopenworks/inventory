// Package backup provides database backup utilities.
package backup

// DB copies the SQLite database at srcPath to destDir.
// It is a placeholder and currently performs no operation.
func DB(_ string, destDir string) error {
	_ = destDir
	return nil
}
