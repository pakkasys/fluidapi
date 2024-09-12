package internal

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/PakkaSys/fluidapi/database/util"
)

const (
	predicateIn = "IN"
	isNotClause = "IS NOT"
	isClause    = "IS"
	sqlNull     = "NULL"
)

// ProcessSelectors processes selectors into where clauses and corresponding
// parameters array.
func ProcessSelectors(selectors []util.Selector) ([]string, []any) {
	var whereColumns []string
	var whereValues []any
	for i := range selectors {
		selector := selectors[i]
		col, vals := processSelector(selector)
		whereColumns = append(whereColumns, col)
		whereValues = append(whereValues, vals...)
	}
	return whereColumns, whereValues
}

func processSelector(
	selector util.Selector,
) (string, []any) {
	switch selector.Predicate {
	case predicateIn:
		return processInSelector(selector)
	default:
		return processDefaultSelector(selector)
	}
}

func processInSelector(
	selector util.Selector,
) (string, []any) {
	value := reflect.ValueOf(selector.Value)
	if value.Kind() == reflect.Slice {
		placeholderCount := value.Len()
		placeholders := createPlaceholders(placeholderCount)
		values := make([]any, placeholderCount)
		for i := 0; i < placeholderCount; i++ {
			values[i] = value.Index(i).Interface()
		}
		column := fmt.Sprintf(
			"`%s`.`%s` %s (%s)",
			selector.Table,
			selector.Field,
			predicateIn,
			placeholders,
		)
		return column, values
	} else {
		// Handle non-slice values as single placeholders
		placeholders := "?"
		values := []any{selector.Value}
		column := fmt.Sprintf(
			"`%s`.`%s` %s (%s)",
			selector.Table,
			selector.Field,
			predicateIn,
			placeholders,
		)
		return column, values
	}
}

func processDefaultSelector(
	selector util.Selector,
) (string, []any) {
	// Check if the value is nil to handle NULL cases
	if selector.Value == nil {
		if selector.Predicate == "=" {
			return fmt.Sprintf(
				"`%s`.`%s` %s %s",
				selector.Table,
				selector.Field,
				isClause,
				sqlNull,
			), nil
		} else if selector.Predicate == "!=" {
			return fmt.Sprintf(
				"`%s`.`%s` %s %s",
				selector.Table,
				selector.Field,
				isNotClause,
				sqlNull,
			), nil
		}
	}

	column := fmt.Sprintf(
		"`%s`.`%s` %s ?",
		selector.Table,
		selector.Field,
		selector.Predicate,
	)
	values := []any{selector.Value}
	return column, values
}

func createPlaceholders(count int) string {
	return strings.TrimSuffix(strings.Repeat("?,", count), ",")
}
