package util

// Selector is a struct that represents a database selector.
type Selector struct {
	Table     string
	Field     string
	Predicate Predicate
	Value     any
}

type Selectors []Selector

// GetByField returns selector with the given field.
//
//   - field: the field to search for
func (s Selectors) GetByField(field string) *Selector {
	for j := range s {
		if s[j].Field == field {
			return &s[j]
		}
	}
	return nil
}

// GetByFields returns selector with the given fields.
//
//   - fields: the fields to search for
func (s Selectors) GetByFields(fields ...string) []Selector {
	selectors := []Selector{}
	for f := range fields {
		for j := range s {
			if s[j].Field == fields[f] {
				selectors = append(selectors, s[j])
			}
		}
	}
	return selectors
}
