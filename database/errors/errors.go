package errors

import (
	"github.com/pakkasys/fluidapi/core/api"
)

const DUPLICATE_ENTRY_ERROR_ID = "DUPLICATE_ENTRY"

func DuplicateEntry(err error) *api.Error {
	return &api.Error{
		ID:   DUPLICATE_ENTRY_ERROR_ID,
		Data: err,
	}
}

const FOREIGN_CONSTRAINT_ERROR_ID = "FOREIGN_CONSTRAINT_ERROR"

func ForeignConstraintError(err error) *api.Error {
	return &api.Error{
		ID:   FOREIGN_CONSTRAINT_ERROR_ID,
		Data: err,
	}
}
