package util

// Predicate represents the predicate of a selector.
type Predicate string

const (
	GREATER          Predicate = ">"
	GREATER_OR_EQUAL Predicate = ">="
	EQUAL            Predicate = "="
	NOT_EQUAL        Predicate = "!="
	LESS             Predicate = "<"
	LESS_OR_EQUAL    Predicate = "<="
	IN               Predicate = "IN"
	NOT_IN           Predicate = "NOT IN"
)
