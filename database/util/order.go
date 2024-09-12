package util

type OrderDirection string

const (
	DIRECTION_ASCENDING  OrderDirection = "ASC"
	DIRECTION_DESCENDING OrderDirection = "DESC"
)

type Order struct {
	Table     string
	Field     string
	Direction OrderDirection
}

func NewOrder(
	table string,
	field string,
	direction OrderDirection,
) *Order {
	return &Order{
		Table:     table,
		Field:     field,
		Direction: direction,
	}
}
