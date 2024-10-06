package order

type OrderDirection string

const (
	DIRECTION_ASC        OrderDirection = "ASC"
	DIRECTION_ASCENDING  OrderDirection = "ASCENDING"
	DIRECTION_DESC       OrderDirection = "DESC"
	DIRECTION_DESCENDING OrderDirection = "DESCENDING"
)

// Directions is a list of all possible order directions.
var Directions []OrderDirection = []OrderDirection{
	DIRECTION_ASC,
	DIRECTION_ASCENDING,
	DIRECTION_DESC,
	DIRECTION_DESCENDING,
}

// Order is used to specify the order of the result set.
type Order struct {
	Field     string         `json:"field"`
	Direction OrderDirection `json:"direction"`
}
