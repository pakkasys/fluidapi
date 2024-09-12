package selector

import (
	"fmt"

	"github.com/pakkasys/fluidapi/endpoint/predicate"
)

type APISelector struct {
	AllowedPredicates []predicate.Predicate
	Validation        string
}

type InputSelector struct {
	Field     string              `json:"field"`
	Predicate predicate.Predicate `json:"predicate"`
	Value     any                 `json:"value"`
}

func (i InputSelector) String() string {
	return fmt.Sprintf("%s %s %s", i.Field, i.Predicate, i.Value)
}

type InputSelectors []InputSelector

// GetByFields returns all selectors with the given fields.
func (i InputSelectors) GetByFields(fields ...string) []InputSelector {
	selectors := InputSelectors{}
	for f := range fields {
		for j := range i {
			if i[j].Field == fields[f] {
				selectors = append(selectors, i[j])
			}
		}
	}
	return selectors
}
