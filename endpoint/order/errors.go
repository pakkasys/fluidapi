package order

import "github.com/pakkasys/fluidapi/core/api"

type InvalidOrderFieldErrorData struct {
	Field string `json:"field"`
}

var InvalidOrderFieldError = api.NewError[InvalidOrderFieldErrorData]("INVALID_ORDER_FIELD")
