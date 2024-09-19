package selector

import (
	"slices"

	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
	"github.com/pakkasys/fluidapi/endpoint/middleware/inputlogic"
	"github.com/pakkasys/fluidapi/endpoint/predicate"
)

type MatchedSelector struct {
	APISelector
	InputSelector
}

type IValidator interface {
	ValidateVariable(fieldName string, obj any, rule string) error
	GetErrorStrings(err error) []string
}

func GetDatabaseSelectorsFromSelectors(
	inputSelectors []InputSelector,
	allowedSelectors map[string]APISelector,
	apiFields map[string]dbfield.DBField,
	validator IValidator,
) ([]util.Selector, error) {
	matchedSelectors, err := MatchAndValidateInputSelectors(
		inputSelectors,
		allowedSelectors,
		validator,
	)
	if err != nil {
		return nil, err
	}

	databaseSelectors, err := ToDatabaseSelectors(apiFields, matchedSelectors)
	if err != nil {
		return nil, err
	}

	return databaseSelectors, nil
}

func MatchAndValidateInputSelectors(
	inputSelectors []InputSelector,
	allowedSelectors map[string]APISelector,
	validator IValidator,
) ([]MatchedSelector, error) {
	var matchedSelectors []MatchedSelector

	for i := range inputSelectors {
		inputSelector := inputSelectors[i]

		// Match the input selector to an allowed selector
		apiSelector, ok := allowedSelectors[inputSelector.Field]
		if !ok {
			return nil, InvalidSelectorFieldError(inputSelector.Field)
		}

		// Validate the input value
		if err := validator.ValidateVariable(
			inputSelector.Field,
			inputSelector.Value,
			apiSelector.Validation,
		); err != nil {
			return nil, inputlogic.ValidationError(
				validator.GetErrorStrings(err),
			)
		}

		matchedSelectors = append(
			matchedSelectors,
			MatchedSelector{
				APISelector:   apiSelector,
				InputSelector: inputSelector,
			},
		)
	}

	return matchedSelectors, nil
}

func ToDatabaseSelectors(
	apiToDatabaseFieldMap map[string]dbfield.DBField,
	matchedSelectors []MatchedSelector,
) ([]util.Selector, error) {
	var databaseSelectors []util.Selector

	for i := range matchedSelectors {
		matchedSelector := matchedSelectors[i]

		// Validate the input predicate
		if !slices.Contains(
			matchedSelector.APISelector.AllowedPredicates,
			matchedSelector.InputSelector.Predicate,
		) {
			return nil, PredicateNotAllowedError(
				matchedSelector.InputSelector.Predicate,
			)
		}

		// Translate the predicate
		translatedPredicate, ok := predicate.
			APIToDatabasePredicates[matchedSelector.InputSelector.Predicate]
		if !ok {
			return nil, InvalidPredicateError(
				matchedSelector.InputSelector.Predicate,
			)
		}

		// Translate the field
		translatedField, ok := apiToDatabaseFieldMap[matchedSelector.Field]
		if !ok {
			return nil, InvalidDatabaseSelectorTranslationError(
				matchedSelector.Field,
			)
		}

		databaseSelectors = append(
			databaseSelectors,
			util.Selector{
				Table:     translatedField.Table,
				Field:     translatedField.Column,
				Predicate: translatedPredicate,
				Value:     matchedSelector.Value,
			},
		)
	}

	return databaseSelectors, nil
}
