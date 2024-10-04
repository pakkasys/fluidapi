package selector

import (
	"fmt"

	"github.com/pakkasys/fluidapi/endpoint/predicate"
)

// Selector represents a data selector.
type Selector struct {
	AllowedPredicates []predicate.Predicate
	Field             string
	Predicate         predicate.Predicate
	Value             any
}

// String returns a string representation of the selector.
func (i Selector) String() string {
	return fmt.Sprintf("%s %s %v", i.Field, i.Predicate, i.Value)
}

type APISelectors []Selector

// GetByFields returns all selectors with the given fields.
//
//   - fields: the fields to search for
func (i APISelectors) GetByFields(fields ...string) []Selector {
	selectors := APISelectors{}
	for f := range fields {
		for j := range i {
			if i[j].Field == fields[f] {
				selectors = append(selectors, i[j])
			}
		}
	}
	return selectors
}
