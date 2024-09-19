package page

type InputPage struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"min=0"`
}

func (p *InputPage) Validate(maxLimit int) error {
	if p.Limit > maxLimit {
		return MaxPageLimitExceeded(maxLimit)
	}
	return nil
}
