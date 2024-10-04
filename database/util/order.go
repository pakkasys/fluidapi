package util

type OrderDirection string

const (
	OrderAsc  OrderDirection = "ASC"
	OrderDesc OrderDirection = "DESC"
)

// Order is used to specify the order of the result set.
type Order struct {
	Table     string
	Field     string
	Direction OrderDirection
}
