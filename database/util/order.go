package util

type OrderDirection string

const (
	DirectionAscending  OrderDirection = "ASC"
	DirectionDescending OrderDirection = "DESC"
)

// Order is used to specify the order of the result set.
type Order struct {
	Table     string
	Field     string
	Direction OrderDirection
}
