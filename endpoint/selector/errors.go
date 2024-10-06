package selector

import (
	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/endpoint/predicate"
)

var NeedAtLeastOneSelectorError = api.NewError[any]("NEED_AT_LEAST_ONE_SELECTOR")

type InvalidDatabaseSelectorTranslationErrorData struct {
	Field string `json:"field"`
}

var InvalidDatabaseSelectorTranslationError = api.NewError[InvalidDatabaseSelectorTranslationErrorData]("INVALID_DATABASE_SELECTOR_TRANSLATION")

type InvalidPredicateErrorData struct {
	Predicate predicate.Predicate `json:""`
}

var InvalidPredicateError = api.NewError[InvalidPredicateErrorData]("INVALID_PREDICATE")

type InvalidSelectorFieldErrorData struct {
	Field string `json:"field"`
}

var InvalidSelectorFieldError = api.NewError[InvalidSelectorFieldErrorData]("INVALID_SELECTOR_FIELD")

type PredicateNotAllowedErrorData struct {
	Predicate predicate.Predicate `json:"predicate"`
}

var PredicateNotAllowedError = api.NewError[PredicateNotAllowedErrorData]("PREDICATE_NOT_ALLOWED")
