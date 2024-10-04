package runner

import (
	"net/http"

	"github.com/pakkasys/fluidapi/database/errors"
	"github.com/pakkasys/fluidapi/endpoint/middleware/inputlogic"
	"github.com/pakkasys/fluidapi/endpoint/order"
	"github.com/pakkasys/fluidapi/endpoint/page"
	"github.com/pakkasys/fluidapi/endpoint/selector"
	"github.com/pakkasys/fluidapi/endpoint/update"
)

var CreateErrors []inputlogic.ExpectedError = []inputlogic.ExpectedError{
	{
		ErrorID:       errors.DuplicateEntryError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: false,
	},
	{
		ErrorID:       errors.ForeignConstraintError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: false,
	},
}

var GetErrors []inputlogic.ExpectedError = []inputlogic.ExpectedError{
	{
		ErrorID:       selector.InvalidPredicateError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.PredicateNotAllowedError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.InvalidSelectorFieldError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       order.InvalidOrderFieldError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       page.MaxPageLimitExceededError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
}

var UpdateErrors []inputlogic.ExpectedError = []inputlogic.ExpectedError{
	{
		ErrorID:       selector.InvalidPredicateError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.InvalidSelectorFieldError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.PredicateNotAllowedError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.NeedAtLeastOneSelectorError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       update.NeedAtLeastOneUpdateError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       order.InvalidOrderFieldError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       errors.DuplicateEntryError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: false,
	},
	{
		ErrorID:       errors.ForeignConstraintError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: false,
	},
}

var DeleteErrors []inputlogic.ExpectedError = []inputlogic.ExpectedError{
	{
		ErrorID:       selector.InvalidPredicateError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.InvalidSelectorFieldError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.PredicateNotAllowedError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.NeedAtLeastOneSelectorError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
}
