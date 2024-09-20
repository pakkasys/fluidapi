package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewMySQLError tests the NewMySQLError function.
func TestNewMySQLError(t *testing.T) {
	number := uint16(1062)
	msg := "Duplicate entry"
	mysqlErr := NewMySQLError(number, msg)

	assert.Equal(t, number, mysqlErr.Number())
	assert.Equal(t, "MySQL Error 1062: Duplicate entry", mysqlErr.Error())
}

// TestMySQLError_Error tests the Error method of MySQLError.
func TestMySQLError_Error(t *testing.T) {
	mysqlErr := NewMySQLError(1452, "Cannot add or update a child row")
	expectedErrMsg := "MySQL Error 1452: Cannot add or update a child row"

	assert.Equal(t, expectedErrMsg, mysqlErr.Error())
}

// TestMySQLError_Is tests the Is method of MySQLError.
func TestMySQLError_Is(t *testing.T) {
	// Case: Error matches
	mysqlErr1 := NewMySQLError(1062, "Duplicate entry")
	mysqlErr2 := NewMySQLError(1062, "Duplicate entry")
	mysqlErr3 := NewMySQLError(1452, "Foreign constraint error")

	assert.True(t, mysqlErr1.Is(mysqlErr2))  // Same error code
	assert.False(t, mysqlErr1.Is(mysqlErr3)) // Different error code

	// Case: Error does not match
	err := errors.New("some other error")
	assert.False(t, mysqlErr1.Is(err))
}

// TestIsMySQLError tests the IsMySQLError function.
func TestIsMySQLError(t *testing.T) {
	mysqlErr := NewMySQLError(
		uint16(MySQLDuplicateEntryErrorCode),
		"Duplicate entry",
	)
	err := errors.New("some other error")

	assert.True(t, IsMySQLError(mysqlErr, MySQLDuplicateEntryErrorCode))
	assert.False(t, IsMySQLError(mysqlErr, MySQLForeignConstraintErrorCode))
	assert.False(t, IsMySQLError(err, MySQLDuplicateEntryErrorCode))
}
