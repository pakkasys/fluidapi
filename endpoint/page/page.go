package page

type InputPage struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"min=0"`
}

func NewInputPage(offset int, limit int) *InputPage {
	return &InputPage{
		Offset: offset,
		Limit:  limit,
	}
}

func ValidatePage(page *InputPage, maxLimit int) (*InputPage, error) {
	if page == nil {
		return NewInputPage(0, maxLimit), nil
	}
	if page.Limit > maxLimit {
		return nil, MaxPageLimitExceeded(maxLimit)
	}

	return page, nil
}
