package page

import "github.com/pakkasys/fluidapi/core/api"

type MaxPageLimitExceededErrorData struct {
	MaxLimit int `json:"max_limit"`
}

var MaxPageLimitExceededError = api.NewError[MaxPageLimitExceededErrorData]("MAX_PAGE_LIMIT_EXCEEDED")

// Page represents a pagination input.
type Page struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"min=0"`
}

// Validate validates the input page.
func (p *Page) Validate(maxLimit int) error {
	if p.Limit > maxLimit {
		return MaxPageLimitExceededError.WithData(
			MaxPageLimitExceededErrorData{
				MaxLimit: maxLimit,
			},
		)
	}
	return nil
}
