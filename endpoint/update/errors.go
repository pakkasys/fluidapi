package update

import "github.com/pakkasys/fluidapi/core/api"

var NeedAtLeastOneUpdateError = api.NewError[any]("NEED_AT_LEAST_ONE_UPDATE")

type InvalidDatabaseUpdateTranslationErrorData struct {
	Field string `json:"field"`
}

var InvalidDatabaseUpdateTranslationError = api.NewError[InvalidDatabaseUpdateTranslationErrorData]("INVALID_DATABASE_UPDATE_TRANSLATION")
