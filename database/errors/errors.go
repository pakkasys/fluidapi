package errors

import (
	"github.com/PakkaSys/fluidapi/core/api"

	"github.com/go-sql-driver/mysql"
)

var DUPLICATE_ENTRY_ERROR_ID = "DUPLICATE_ENTRY"

func DuplicateEntry(mySQLError *mysql.MySQLError) *api.Error {
	return &api.Error{
		ID:   DUPLICATE_ENTRY_ERROR_ID,
		Data: mySQLError,
	}

}

var FOREIGN_CONSTRAINT_ERROR_ID = "FOREIGN_CONSTRAINT_ERROR"

func ForeignConstraintError(mySQLError *mysql.MySQLError) *api.Error {
	return &api.Error{
		ID:   FOREIGN_CONSTRAINT_ERROR_ID,
		Data: mySQLError,
	}
}
