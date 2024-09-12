package update

import "github.com/PakkaSys/fluidapi/core/api"

type InvalidDatabaseUpdateTranslationErrorData struct {
	Field string `json:"field"`
}

var INVALID_DATABASE_UPDATE_TRANSLATION_ERROR_ID = "INVALID_DATABASE_UPDATE_TRANSLATION"

func InvalidDatabaseUpdateTranslationError(field string) *api.Error {
	return &api.Error{
		ID: INVALID_DATABASE_UPDATE_TRANSLATION_ERROR_ID,
		Data: InvalidDatabaseUpdateTranslationErrorData{
			Field: field,
		},
	}
}

type InvalidUpdateFieldErrorData struct {
	Field string `json:"field"`
}

var INVALID_UPDATE_FIELD_ERROR_ID = "INVALID_UPDATE_FIELD"

func InvalidUpdateFieldError(field string) *api.Error {
	return &api.Error{
		ID: INVALID_UPDATE_FIELD_ERROR_ID,
		Data: InvalidUpdateFieldErrorData{
			Field: field,
		},
	}
}

var NEED_AT_LEAST_ONE_UPDATE_ERROR_ID = "NEED_AT_LEAST_ONE_UPDATE"

func NeedAtLeastOneUpdateError() *api.Error {
	return &api.Error{
		ID: NEED_AT_LEAST_ONE_UPDATE_ERROR_ID,
	}
}
