package selector

import (
	"slices"

	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
	"github.com/pakkasys/fluidapi/endpoint/predicate"
)

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
		dbPredicate, ok := predicate.
			APIToDatabasePredicates[selector.Predicate]
		if !ok {
			return nil, InvalidPredicateError.WithData(
				InvalidPredicateErrorData{Predicate: selector.Predicate},
			)
		}

		// Translate the field
		dbField, ok := apiToDBFieldMap[selector.Field]
		if !ok {
			return nil, InvalidDatabaseSelectorTranslationError.WithData(
				InvalidDatabaseSelectorTranslationErrorData{Field: selector.Field},
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
