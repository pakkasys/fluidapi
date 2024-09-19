package internal

import (
	"errors"
	"fmt"
)

// DBError is an interface representing a database error with a number/code.
type DBError interface {
	Number() uint16
	Error() string
	Is(err error) bool
}

// MySQLErrorCode represents MySQL error codes.
type MySQLErrorCode uint16

const (
	MySQLForeignConstraintErrorCode MySQLErrorCode = 1452
	MySQLDuplicateEntryErrorCode    MySQLErrorCode = 1062
)

// MySQLError represents a simplified version of a MySQL error.
type MySQLError struct {
	number uint16
	msg    string
}

// NewMySQLError creates a new MySQLError with a given number and message.
func NewMySQLError(number uint16, msg string) *MySQLError {
	return &MySQLError{number: number, msg: msg}
}

// Number returns the error number/code.
func (e *MySQLError) Number() uint16 {
	return e.number
}

// Error implements the error interface.
func (e *MySQLError) Error() string {
	return fmt.Sprintf("MySQL Error %d: %s", e.number, e.msg)
}

// Is allows comparison of MySQLError objects.
func (e *MySQLError) Is(err error) bool {
	if merr, ok := err.(*MySQLError); ok {
		return merr.number == e.number
	}
	return false
}

// IsMySQLError checks if an error is a MySQLError and matches the given code.
func IsMySQLError(err error, code MySQLErrorCode) bool {
	var mysqlErr *MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number() == uint16(code)
	}
	return false
}
