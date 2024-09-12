package predicate

type Predicate string

const (
	GREATER       Predicate = ">"
	GREATER_SHORT Predicate = "GT"

	GREATER_OR_EQUAL       Predicate = ">="
	GREATER_OR_EQUAL_SHORT Predicate = "GE"

	EQUAL       Predicate = "="
	EQUAL_SHORT Predicate = "EQ"

	NOT_EQUAL       Predicate = "!="
	NOT_EQUAL_SHORT Predicate = "NE"

	LESS       Predicate = "<"
	LESS_SHORT Predicate = "LT"

	LESS_OR_EQUAL       Predicate = "<="
	LESS_OR_EQUAL_SHORT Predicate = "LE"

	IN     Predicate = "IN"
	NOT_IN Predicate = "NOT_IN"
)

var AllPredicates = []Predicate{
	GREATER,
	GREATER_SHORT,
	GREATER_OR_EQUAL,
	GREATER_OR_EQUAL_SHORT,
	EQUAL,
	EQUAL_SHORT,
	NOT_EQUAL,
	NOT_EQUAL_SHORT,
	LESS,
	LESS_SHORT,
	LESS_OR_EQUAL,
	LESS_OR_EQUAL_SHORT,
	IN,
	NOT_IN,
}

var OnlyEqualPredicate = []Predicate{
	EQUAL,
	EQUAL_SHORT,
}

var EqualAndNotEqualPredicates = []Predicate{
	EQUAL,
	EQUAL_SHORT,
	NOT_EQUAL,
	NOT_EQUAL_SHORT,
}

var OnlyGreaterPredicates = []Predicate{
	GREATER_OR_EQUAL,
	GREATER_OR_EQUAL_SHORT,
	GREATER,
	GREATER_SHORT,
}

var OnlyLessPredicates = []Predicate{
	LESS_OR_EQUAL,
	LESS_OR_EQUAL_SHORT,
	LESS,
	LESS_SHORT,
}

var OnlyInAndNotInPredicates = []Predicate{
	IN,
	NOT_IN,
}
