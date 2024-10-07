package selector

import (
	"fmt"
	"slices"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
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

// Selector represents a data selector that specifies criteria for filtering
// data based on fields, predicates, and values.
type Selector struct {
	// Predicates allowed for this selector
	AllowedPredicates []predicate.Predicate
	// The name of the field being filtered
	Field string
	// The predicate for filtering
	Predicate predicate.Predicate
	// The value to filter by
	Value any
}

// Selectors represents a collection of selectors used for filtering data.
type Selectors []Selector

// String returns a string representation of the selector.
// It is useful for debugging and logging purposes.
//
// Returns:
// - A formatted string showing the field, predicate, and value.
func (i Selector) String() string {
	return fmt.Sprintf("%s %s %v", i.Field, i.Predicate, i.Value)
}

// GetByFields returns all selectors that match the given fields.
//
// Parameters:
// - fields: The fields to search for in the selectors.
//
// Returns:
// - A slice of selectors that match the provided field names.
func (i Selectors) GetByFields(fields ...string) []Selector {
	selectors := Selectors{}
	for f := range fields {
		for j := range i {
			if i[j].Field == fields[f] {
				selectors = append(selectors, i[j])
			}
		}
	}
	return selectors
}

// ToDBSelectors converts a slice of API-level selectors to database selectors.
// It validates predicates and translates the fields and predicates for use with
// the database.
//
// Parameters:
//   - selectors: A slice of API-level selectors that specify the criteria for
//     selecting data.
//   - apiToDBFieldMap: A map translating API field names to their corresponding
//     database field definitions.
//
// Returns:
//   - A slice of util.Selector, which represents the translated database
//     selectors.
//   - An error if any validation fails, such as invalid predicates or unknown
//     fields.
func ToDBSelectors(
	selectors []Selector,
	apiToDBFieldMap map[string]dbfield.DBField,
) ([]util.Selector, error) {
	var databaseSelectors []util.Selector

	for i := range selectors {
		selector := selectors[i]

		// Validate the input predicate
		if !slices.Contains(
			selector.AllowedPredicates,
			selector.Predicate,
		) {
			return nil, PredicateNotAllowedError.WithData(
				PredicateNotAllowedErrorData{Predicate: selector.Predicate},
			)
		}

		// Translate the predicate
		dbPredicate, ok := predicate.ToDBPredicates[selector.Predicate]
		if !ok {
			return nil, InvalidPredicateError.WithData(
				InvalidPredicateErrorData{Predicate: selector.Predicate},
			)
		}

		// Translate the field
		dbField, ok := apiToDBFieldMap[selector.Field]
		if !ok {
			return nil, InvalidDatabaseSelectorTranslationError.WithData(
				InvalidDatabaseSelectorTranslationErrorData{
					Field: selector.Field,
				},
			)
		}

		databaseSelectors = append(
			databaseSelectors,
			util.Selector{
				Table:     dbField.Table,
				Field:     dbField.Column,
				Predicate: dbPredicate,
				Value:     selector.Value,
			},
		)
	}

	return databaseSelectors, nil
}
