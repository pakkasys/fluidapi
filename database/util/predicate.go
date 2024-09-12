package util

type Predicate string

func (predicate Predicate) String() string {
	return string(predicate)
}

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
