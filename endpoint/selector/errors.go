package selector

import (
	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/endpoint/predicate"
)

type InvalidDatabaseSelectorTranslationErrorData struct {
	Field string `json:"field"`
}

var INVALID_DATABASE_SELECTOR_TRANSLATION_ERROR_ID = "INVALID_DATABASE_SELECTOR_TRANSLATION"

func InvalidDatabaseSelectorTranslationError(field string) *api.Error {
	return &api.Error{
		ID: INVALID_DATABASE_SELECTOR_TRANSLATION_ERROR_ID,
		Data: InvalidDatabaseSelectorTranslationErrorData{
			Field: field,
		},
	}
}

type InvalidPredicateErrorData struct {
	Predicate predicate.Predicate `json:""`
}

var INVALID_PREDICATE_ERROR_ID = "INVALID_PREDICATE"

func InvalidPredicateError(predicate predicate.Predicate) *api.Error {
	return &api.Error{
		ID: INVALID_PREDICATE_ERROR_ID,
		Data: InvalidPredicateErrorData{
			Predicate: predicate,
		},
	}
}

type InvalidSelectorFieldErrorData struct {
	Field string `json:"field"`
}

var INVALID_SELECTOR_FIELD_ERROR_ID = "INVALID_SELECTOR_FIELD"

func InvalidSelectorFieldError(field string) *api.Error {
	return &api.Error{
		ID: INVALID_SELECTOR_FIELD_ERROR_ID,
		Data: InvalidSelectorFieldErrorData{
			Field: field,
		},
	}
}

var NEED_AT_LEAST_ONE_SELECTOR_ERROR_ID = "NEED_AT_LEAST_ONE_SELECTOR"

func NeedAtLeastOneSelectorError() *api.Error {
	return &api.Error{
		ID: NEED_AT_LEAST_ONE_SELECTOR_ERROR_ID,
	}
}

type PredicateNotAllowedErrorData struct {
	Predicate predicate.Predicate `json:""`
}

var PREDICATE_NOT_ALLOWED_ERROR_ID = "PREDICATE_NOT_ALLOWED"

func PredicateNotAllowedError(
	predicate predicate.Predicate,
) *api.Error {
	return &api.Error{
		ID: PREDICATE_NOT_ALLOWED_ERROR_ID,
		Data: PredicateNotAllowedErrorData{
			Predicate: predicate,
		},
	}
}
