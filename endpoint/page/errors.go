package page

import "github.com/pakkasys/fluidapi/core/api"

type MaxPageLimitExceededErrorData struct {
	MaxLimit int `json:"max_limit"`
}

var MaxPageLimitExceededError = api.NewError[MaxPageLimitExceededErrorData]("MAX_PAGE_LIMIT_EXCEEDED")
