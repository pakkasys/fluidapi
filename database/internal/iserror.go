package internal

import "github.com/go-sql-driver/mysql"

// IsForeignConstraintError checks if an error is a foreign constraint error.
func IsForeignConstraintError(err *mysql.MySQLError) bool {
	return err.Number == 1452
}

// IsDuplicateEntryError checks if an error is a duplicate entry error.
func IsDuplicateEntryError(err *mysql.MySQLError) bool {
	return err.Number == 1062
}
