package page

import "github.com/pakkasys/fluidapi/core/api"

type MaxPageLimitExceededErrorData struct {
	MaxLimit int `json:"max_limit"`
}

var MAX_PAGE_LIMIT_EXCEEDED_ERROR_ID = "MAX_PAGE_LIMIT_EXCEEDED"

func MaxPageLimitExceeded(maxLimit int) *api.Error {
	return &api.Error{
		ID: MAX_PAGE_LIMIT_EXCEEDED_ERROR_ID,
		Data: MaxPageLimitExceededErrorData{
			MaxLimit: maxLimit,
		},
	}
}
