package runner

import (
	"net/http"

	"github.com/PakkaSys/fluidapi/database/errors"
	"github.com/PakkaSys/fluidapi/endpoint/middleware/inputlogic"
	"github.com/PakkaSys/fluidapi/endpoint/order"
	"github.com/PakkaSys/fluidapi/endpoint/page"
	"github.com/PakkaSys/fluidapi/endpoint/selector"
	"github.com/PakkaSys/fluidapi/endpoint/update"
)

var ExpectedErrorsCreate []inputlogic.ExpectedError = []inputlogic.ExpectedError{
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

var ExpectedErrorsGet []inputlogic.ExpectedError = []inputlogic.ExpectedError{
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

var ExpectedErrorsUpdate []inputlogic.ExpectedError = []inputlogic.ExpectedError{
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

var ExpectedErrorsDelete []inputlogic.ExpectedError = []inputlogic.ExpectedError{
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
