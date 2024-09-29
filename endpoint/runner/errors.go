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
		ErrorID:       errors.DUPLICATE_ENTRY_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: false,
	},
	{
		ErrorID:       errors.FOREIGN_CONSTRAINT_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: false,
	},
}

var GetErrors []inputlogic.ExpectedError = []inputlogic.ExpectedError{
	{
		ErrorID:       selector.INVALID_PREDICATE_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.PREDICATE_NOT_ALLOWED_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.INVALID_SELECTOR_FIELD_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       order.INVALID_ORDER_FIELD_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       page.MAX_PAGE_LIMIT_EXCEEDED_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
}

var UpdateErrors []inputlogic.ExpectedError = []inputlogic.ExpectedError{
	{
		ErrorID:       selector.INVALID_PREDICATE_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.INVALID_SELECTOR_FIELD_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.PREDICATE_NOT_ALLOWED_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.NEED_AT_LEAST_ONE_SELECTOR_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       update.NEED_AT_LEAST_ONE_UPDATE_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       update.INVALID_UPDATE_FIELD_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       order.INVALID_ORDER_FIELD_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       errors.DUPLICATE_ENTRY_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: false,
	},
	{
		ErrorID:       errors.FOREIGN_CONSTRAINT_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: false,
	},
}

var DeleteErrors []inputlogic.ExpectedError = []inputlogic.ExpectedError{
	{
		ErrorID:       selector.INVALID_PREDICATE_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.INVALID_SELECTOR_FIELD_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.PREDICATE_NOT_ALLOWED_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
	{
		ErrorID:       selector.NEED_AT_LEAST_ONE_SELECTOR_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
}
