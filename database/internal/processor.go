package internal

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pakkasys/fluidapi/database/util"
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
	for _, selector := range selectors {
		col, vals := processSelector(selector)
		whereColumns = append(whereColumns, col)
		whereValues = append(whereValues, vals...)
	}
	return whereColumns, whereValues
}

func processSelector(selector util.Selector) (string, []any) {
	if selector.Predicate == predicateIn {
		return processInSelector(selector)
	}
	return processDefaultSelector(selector)
}

func processInSelector(selector util.Selector) (string, []any) {
	value := reflect.ValueOf(selector.Value)
	if value.Kind() == reflect.Slice {
		placeholders, values := createPlaceholdersAndValues(value)
		column := fmt.Sprintf(
			"`%s`.`%s` %s (%s)",
			selector.Table,
			selector.Field,
			predicateIn,
			placeholders,
		)
		return column, values
	}
	// If value is not a slice, treat as a single value
	return fmt.Sprintf(
		"`%s`.`%s` %s (?)",
		selector.Table,
		selector.Field,
		predicateIn,
	), []any{selector.Value}
}

func processDefaultSelector(selector util.Selector) (string, []any) {
	if selector.Value == nil {
		return processNullSelector(selector)
	}
	if selector.Table == "" {
		return fmt.Sprintf(
			"`%s` %s ?",
			selector.Field,
			selector.Predicate,
		), []any{selector.Value}
	} else {
		return fmt.Sprintf(
			"`%s`.`%s` %s ?",
			selector.Table,
			selector.Field,
			selector.Predicate,
		), []any{selector.Value}
	}
}

func processNullSelector(selector util.Selector) (string, []any) {
	if selector.Predicate == "=" {
		return buildNullClause(selector, isClause), nil
	}
	if selector.Predicate == "!=" {
		return buildNullClause(selector, isNotClause), nil
	}
	return "", nil
}

func buildNullClause(selector util.Selector, clause string) string {
	if selector.Table == "" {
		return fmt.Sprintf("`%s` %s %s", selector.Field, clause, sqlNull)
	}
	return fmt.Sprintf(
		"`%s`.`%s` %s %s",
		selector.Table,
		selector.Field,
		clause,
		sqlNull,
	)
}

func createPlaceholdersAndValues(value reflect.Value) (string, []any) {
	placeholderCount := value.Len()
	placeholders := createPlaceholders(placeholderCount)
	values := make([]any, placeholderCount)
	for i := 0; i < placeholderCount; i++ {
		values[i] = value.Index(i).Interface()
	}
	return placeholders, values
}

func createPlaceholders(count int) string {
	return strings.TrimSuffix(strings.Repeat("?,", count), ",")
}
