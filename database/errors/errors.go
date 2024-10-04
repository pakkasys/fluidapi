package errors

import (
	"github.com/pakkasys/fluidapi/core/api"
)

var (
	DuplicateEntryError    = api.NewError[error]("DUPLICATE_ENTRY")
	ForeignConstraintError = api.NewError[error]("FOREIGN_CONSTRAINT_ERROR")
)
