package util

type Selector struct {
	Table     string
	Field     string
	Predicate Predicate
	Value     any
}

func NewSelector(
	table string,
	field string,
	predicate Predicate,
	value any,
) *Selector {
	return &Selector{
		Table:     table,
		Field:     field,
		Predicate: predicate,
		Value:     value,
	}
}

type Selectors []Selector

// GetByField returns selector with the given fields.
func (s Selectors) GetByField(field string) *Selector {
	for j := range s {
		if s[j].Field == field {
			return &s[j]
		}
	}
	return nil
}

// GetByFields returns selector with the given fields.
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
